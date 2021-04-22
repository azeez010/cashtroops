package ops

import (
	"github.com/adigunhammedolalekan/cashtroops/errors"
	"github.com/adigunhammedolalekan/cashtroops/libs/paystackclient"
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AccountOps interface {
	AddBeneficiary(userId string, beneficiary *types.CreateBeneficiaryOpts) (*types.Beneficiary, error)
	RemoveBeneficiary(owner, beneficiaryId string) error
	ListBeneficiariesFor(userId string) ([]*types.Beneficiary, error)
	GetBeneficiaryByAttr(attr string, value interface{}) (*types.Beneficiary, error)
	SetBanks(banks []types.Bank)
	GetBanks() []types.Bank
	ResolveAccount(accountNumber, bankCode string) (*types.BankAccount, error)
}

type accountOps struct {
	db     *gorm.DB
	logger *logrus.Logger
	ps     paystackclient.Client
	banks  []types.Bank
}

func NewAccountOps(db *gorm.DB, ps paystackclient.Client, logger *logrus.Logger) AccountOps {
	return &accountOps{
		db:     db,
		ps:     ps,
		logger: logger,
	}
}

func (a *accountOps) SetBanks(banks []types.Bank) {
	a.banks = banks
}

func (a *accountOps) AddBeneficiary(userId string, beneficiary *types.CreateBeneficiaryOpts) (*types.Beneficiary, error) {
	existing, err := a.GetBeneficiaryByAttr("account_number", beneficiary.AccountNumber)
	if err == nil && existing.AccountNumber != "" && existing.Owner == userId {
		return nil, errors.New(http.StatusConflict, "beneficiary with this account number has already been added")
	}
	bf := types.NewBeneficiary(beneficiary, userId)
	if err := a.db.Table("beneficiaries").Create(bf).Error; err != nil {
		a.logger.WithError(err).Error("failed to create beneficiary")
		return nil, err
	}
	return bf, nil
}

func (a *accountOps) RemoveBeneficiary(owner, id string) error {
	bf, err := a.GetBeneficiaryByAttr("id", id)
	if err != nil {
		return err
	}
	if bf.Owner != owner {
		return errors.New(http.StatusForbidden, "You cannot remove beneficiary that does not belong to you")
	}
	return a.db.Table("beneficiaries").Where("id = ?", id).Delete(&types.Beneficiary{}).Error
}

func (a *accountOps) ListBeneficiariesFor(userId string) ([]*types.Beneficiary, error) {
	values := make([]*types.Beneficiary, 0)
	err := a.db.Table("beneficiaries").Where("owner = ? AND hidden = ?", userId, false).Find(&values).Error
	return values, err
}

func (a *accountOps) GetBeneficiaryByAttr(attr string, value interface{}) (*types.Beneficiary, error) {
	bf := &types.Beneficiary{}
	err := a.db.Table("beneficiaries").Where(attr+" = ?", value).First(bf).Error
	return bf, err
}

func (a *accountOps) GetBanks() []types.Bank {
	return a.banks
}

func (a *accountOps) ResolveAccount(accountNumber, bankCode string) (*types.BankAccount, error) {
	return a.ps.ResolveAccountNumber(accountNumber, bankCode)
}
