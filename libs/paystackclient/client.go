package paystackclient

import (
	"fmt"
	"github.com/adigunhammedolalekan/cashtroops/libs"
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	baseUrl = "https://api.paystack.co/"
)

type TransferRecipientBody struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	Currency      string `json:"currency"`
}

type TransferRecipient struct {
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"createdAt"`
	Currency      string    `json:"currency"`
	Domain        string    `json:"domain"`
	ID            int       `json:"id"`
	Integration   int       `json:"integration"`
	Name          string    `json:"name"`
	RecipientCode string    `json:"recipient_code"`
	Type          string    `json:"type"`
	UpdatedAt     time.Time `json:"updatedAt"`
	IsDeleted     bool      `json:"is_deleted"`
	Details       struct {
		AuthorizationCode interface{} `json:"authorization_code"`
		AccountNumber     string      `json:"account_number"`
		AccountName       interface{} `json:"account_name"`
		BankCode          string      `json:"bank_code"`
		BankName          string      `json:"bank_name"`
	} `json:"details"`
}

type InitiateTransferRequest struct {
	Source    string `json:"source"`
	Amount    int64  `json:"amount"`
	Recipient string `json:"recipient"`
	Reason    string `json:"reason"`
}

type Client interface {
	CreateTransferRecipient(recipient *TransferRecipientBody) (*TransferRecipient, error)
	InitiateTransfer(req *InitiateTransferRequest) (*types.Transfer, error)
	ResolveAccountNumber(accountNumber, bankCode string) (*types.BankAccount, error)
}

type paystackClient struct {
	httpClient *libs.HttpClient
	logger     *logrus.Logger
}

func New(bearerToken string, logger *logrus.Logger) Client {
	return &paystackClient{httpClient: libs.NewHttpClient(logger, bearerToken)}
}

func (ps *paystackClient) CreateTransferRecipient(recipient *TransferRecipientBody) (*TransferRecipient, error) {
	var data struct {
		Data *TransferRecipient `json:"data"`
	}
	u := fmt.Sprintf("%s/transferrecipient", baseUrl)
	err := ps.httpClient.Do(u, "POST", recipient, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (ps *paystackClient) InitiateTransfer(req *InitiateTransferRequest) (*types.Transfer, error) {
	var data struct {
		Data *types.Transfer `json:"data"`
	}
	u := fmt.Sprintf("%s/transfer", baseUrl)
	err := ps.httpClient.Do(u, "POST", req, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (ps *paystackClient) ResolveAccountNumber(accountNumber, bankCode string) (*types.BankAccount, error) {
	account := &types.BankAccount{}
	u := fmt.Sprintf("%s/bank/resolve?account_number=%s&bank_code=%s", baseUrl, accountNumber, bankCode)
	err := ps.httpClient.Do(u, "GET", nil, account)
	return account, err
}
