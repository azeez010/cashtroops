package types

import (
	"encoding/json"
	"time"
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

type Bank struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Code     string `json:"code"`
	Currency string `json:"currency"`
	ID       int64  `json:"id"`
	Longcode string `json:"longcode"`
}

type BankAccount struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	BankId        int64  `json:"bank_id"`
}

type TransferEvent struct {
	Amount      int         `json:"amount"`
	Currency    string      `json:"currency"`
	Domain      string      `json:"domain"`
	Failures    interface{} `json:"failures"`
	ID          int         `json:"id"`
	Integration struct {
		ID           int    `json:"id"`
		IsLive       bool   `json:"is_live"`
		BusinessName string `json:"business_name"`
	} `json:"integration"`
	Reason        string      `json:"reason"`
	Reference     string      `json:"reference"`
	Source        string      `json:"source"`
	SourceDetails interface{} `json:"source_details"`
	Status        string      `json:"status"`
	TitanCode     interface{} `json:"titan_code"`
	TransferCode  string      `json:"transfer_code"`
	TransferredAt interface{} `json:"transferred_at"`
	Recipient     struct {
		Active        bool        `json:"active"`
		Currency      string      `json:"currency"`
		Description   string      `json:"description"`
		Domain        string      `json:"domain"`
		Email         interface{} `json:"email"`
		ID            int         `json:"id"`
		Integration   int         `json:"integration"`
		Metadata      interface{} `json:"metadata"`
		Name          string      `json:"name"`
		RecipientCode string      `json:"recipient_code"`
		Type          string      `json:"type"`
		IsDeleted     bool        `json:"is_deleted"`
		Details       struct {
			AccountNumber string      `json:"account_number"`
			AccountName   interface{} `json:"account_name"`
			BankCode      string      `json:"bank_code"`
			BankName      string      `json:"bank_name"`
		} `json:"details"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"recipient"`
	Session struct {
		Provider interface{} `json:"provider"`
		ID       interface{} `json:"id"`
	} `json:"session"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransferRecipient struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	Currency      string `json:"currency"`
}

func (req *InitPaymentRequest) AmountInt() int64 {
	if value, err := req.Amount.Int64(); err == nil {
		return value
	}
	return 0
}
