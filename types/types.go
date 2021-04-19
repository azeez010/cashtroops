package types

import (
	"encoding/json"
)

type CreateUserOpts struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type CreateBeneficiaryOpts struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	Hidden        bool   `json:"hidden"`
}

type MailRequest struct {
	User  string
	Email string
	Title string
	Body  string
}

type InitPaymentResponse struct {
	AddressUsed string `json:"address_used"`
	Coin        string `json:"coin"`
	PaymentId   string `json:"payment_id"`
}

type PaymentBeneficiary struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
}

type InitPaymentRequest struct {
	BeneficiaryId string              `json:"beneficiary_id"`
	Beneficiary   *PaymentBeneficiary `json:"beneficiary"`
	Amount        json.Number         `json:"amount"`
	Coin          string              `json:"coin"`
}

func (req *InitPaymentRequest) AmountInt() int64 {
	if value, err := req.Amount.Int64(); err == nil {
		return value
	}
	return 0
}
