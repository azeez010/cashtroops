package ops

import (
	"github.com/adigunhammedolalekan/cashtroops/errors"
	"github.com/adigunhammedolalekan/cashtroops/libs/bc"
	"github.com/adigunhammedolalekan/cashtroops/libs/priceclient"
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/btcsuite/btcutil"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type PaymentOps interface {
	InitializePayment(userId string, req *types.InitPaymentRequest) (*types.InitPaymentResponse, error)
	VerifyWebHookId(webHookId string) error
	FinalizePayment(address string, amount int64) error
	GetPaymentByAddress(address string) (*types.Payment, error)
	GetCurrentRate(currency string) (*types.Rate, error)
	ProcessPayment(payment *types.Payment) error
	ListPayments(userId string) ([]*types.Payment, error)
	InitRate(pair string, value int64) error
}

type paymentOps struct {
	db          *gorm.DB
	userOps     UserOps
	accountOps  AccountOps
	bcClient    bc.Client
	priceClient priceclient.Client
	logger      *logrus.Logger
}

func NewPaymentOps(
	db *gorm.DB,
	bcClient bc.Client,
	ops UserOps,
	accountOps AccountOps, logger *logrus.Logger) PaymentOps {
	return &paymentOps{
		db:          db,
		bcClient:    bcClient,
		priceClient: priceclient.New(logger),
		userOps:     ops,
		accountOps:  accountOps,
		logger:      logger,
	}
}

func (p *paymentOps) InitializePayment(userId string, req *types.InitPaymentRequest) (*types.InitPaymentResponse, error) {
	beneficiaryId := req.BeneficiaryId
	if beneficiaryId == "" && req.Beneficiary != nil {
		newBeneficiary := types.NewBeneficiary(&types.CreateBeneficiaryOpts{
			AccountName:   req.Beneficiary.AccountName,
			AccountNumber: req.Beneficiary.AccountNumber,
			BankName:      req.Beneficiary.BankName,
			BankCode:      req.Beneficiary.BankCode,
			Hidden:        true,
		}, userId)
		if err := p.db.Table("beneficiaries").Create(newBeneficiary).Error; err != nil {
			p.logger.WithError(err).Error("failed to create hidden beneficiary")
			return nil, errors.New(http.StatusInternalServerError, "failed to process transaction at this time. please retry later.")
		}
		beneficiaryId = newBeneficiary.ID.String()
	}
	payment := &types.Payment{
		UserId:        userId,
		Amount:        req.AmountInt(),
		Coin:          req.Coin,
		BeneficiaryId: beneficiaryId,
		Ts:            time.Now(),
	}
	tx := p.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	if err := tx.Table("payments").Create(payment).Error; err != nil {
		tx.Rollback()
		p.logger.WithError(err).Error("failed to create payment body")
		return nil, errors.New(http.StatusInternalServerError, "failed to process transaction at this time. please retry later.")
	}
	addr, err := p.bcClient.GenerateAddress()
	if err != nil {
		tx.Rollback()
		return nil, errors.New(http.StatusInternalServerError, "failed to process transaction at this time. please retry later.")
	}
	newAddress := &types.Address{
		Public:   addr.Public,
		Private:  addr.Private,
		Provider: "BLOCKCYPHER",
		Coin:     req.Coin,
		UserId:   userId,
		Ts:       time.Now(),
	}
	if err := tx.Table("addresses").Create(newAddress).Error; err != nil {
		tx.Rollback()
		p.logger.WithError(err).Error("failed to create generated address")
		return nil, errors.New(http.StatusInternalServerError, "failed to process transaction at this time. please retry later.")
	}
	err = tx.Table("payments").Where("id = ?", payment.ID.String()).UpdateColumn("address_used", addr.Address).Error
	if err != nil {
		tx.Rollback()
		p.logger.WithError(err).Error("failed to update payment crypto address")
		return nil, errors.New(http.StatusInternalServerError, "failed to process transaction at this time. please retry later.")
	}
	event := bc.Event{
		Event:   bc.EventTypeTxConfirmation,
		URL:     "https://webhook.site/59ae6da0-029a-4368-9696-5f533b3c36a7",
		Address: addr.Address,
	}
	newEvent, err := p.bcClient.AddHook(event)
	if err != nil {
		tx.Rollback()
		p.logger.WithError(err).Error("failed to setup notification for address")
		return nil, errors.New(http.StatusInternalServerError, "failed to process transaction at this time. please retry later.")
	}
	hook := &types.Hook{
		URL:     newEvent.URL,
		HookId:  newEvent.ID,
		Address: newEvent.Address,
		Ts:      time.Now(),
	}
	if err := tx.Table("hooks").Create(hook).Error; err != nil {
		tx.Rollback()
		p.logger.WithError(err).Error("failed to create hook")
		return nil, errors.New(http.StatusInternalServerError, "failed to process transaction at this time. please retry later.")
	}
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(http.StatusInternalServerError, "failed to process transaction at this time. please retry later.")
	}
	return &types.InitPaymentResponse{
		AddressUsed: addr.Address,
		Coin:        req.Coin,
		PaymentId:   payment.ID.String(),
	}, nil
}

func (p *paymentOps) VerifyWebHookId(webHookId string) error {
	hooks, err := p.bcClient.ListHooks()
	if err != nil {
		return err
	}
	for _, next := range hooks {
		if next.ID == webHookId {
			return nil
		}
	}
	return errors.New(http.StatusUnauthorized, "webHook ID is not found")
}

func (p *paymentOps) FinalizePayment(address string, amount int64) error {
	payment, err := p.GetPaymentByAddress(address)
	if err != nil {
		p.logger.WithError(err).Error("payment not found")
		return err
	}
	if payment.Status == types.DONE {
		p.logger.WithFields(logrus.Fields{
			"status": payment.Status,
			"id":     payment.ID.String(),
		}).Info("payment has already been processed")
		return errors.New(http.StatusConflict, "payment has already been processed")
	}
	currentBtcPrice, err := p.priceClient.CurrentPrice()
	if err != nil {
		p.logger.WithError(err).Error("failed to get current BTC price")
		return errors.New(http.StatusInternalServerError, "failed to finalize transaction")
	}
	rate, err := p.GetCurrentRate("USD-NGN")
	if err != nil {
		p.logger.WithError(err).Error("failed to get current rate")
		return errors.New(http.StatusInternalServerError, "failed to get USD-NGN rate")
	}
	// convert amount to BTC, then USD, then NGN
	btcAmount := btcutil.Amount(amount).ToUnit(btcutil.AmountBTC)
	amountInUsd := btcAmount * currentBtcPrice.FloatAmount()
	amountInKobo := (int64(amountInUsd) * rate.Value) * 100
	payment.KoboAmount = amountInKobo
	payment.UsdAmount = amountInUsd
	payment.BtcAmount = btcAmount
	err = p.db.Table("payments").Where("id = ?", payment.ID.String()).Update(payment).Error
	if err != nil {
		return err
	}
	return p.ProcessPayment(payment)
}

func (p *paymentOps) GetPaymentByAddress(address string) (*types.Payment, error) {
	payment := &types.Payment{}
	err := p.db.Table("payments").Where("address_used = ?", address).First(payment).Error
	return payment, err
}

func (p *paymentOps) GetCurrentRate(currencyPair string) (*types.Rate, error) {
	rate := &types.Rate{}
	err := p.db.Table("rates").Where("currency_pair = ?", currencyPair).First(rate).Error
	return rate, err
}

func (p *paymentOps) ProcessPayment(payment *types.Payment) error {
	return nil
}

func (p *paymentOps) ListPayments(userId string) ([]*types.Payment, error) {
	values := make([]*types.Payment, 0)
	err := p.db.Table("payments").Where("user_id = ?", userId).Find(&values).Error
	return values, err
}

func (p *paymentOps) InitRate(pair string, value int64) error {
	rate := &types.Rate{
		CurrencyPair: pair,
		Value:        value,
		Ts:           time.Now(),
	}
	if value, err := p.GetCurrentRate(rate.CurrencyPair); err == nil && value.ID.String() != "" {
		p.db.Table("rates").Where("currency_pair = ?", rate.CurrencyPair).Delete(&types.Rate{})
	}
	return p.db.Table("rates").Create(rate).Error
}
