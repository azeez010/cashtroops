package types

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Beneficiary struct {
	ID                  uuid.UUID `json:"id" gorm:"primary_key"`
	AccountName         string    `json:"account_name"`
	AccountNumber       string    `json:"account_number"`
	BankName            string    `json:"bank_name"`
	BankCode            string    `json:"bank_code"`
	TransferRecipientId string    `json:"transfer_recipient_id"`
	Owner               string    `json:"owner"`
	Hidden              bool      `json:"hidden"`
}

func NewBeneficiary(opts *CreateBeneficiaryOpts, owner string) *Beneficiary {
	return &Beneficiary{
		ID:            uuid.New(),
		AccountName:   opts.AccountName,
		AccountNumber: opts.AccountNumber,
		BankName:      opts.BankName,
		Owner:         owner,
		BankCode:      opts.BankCode,
	}
}

func (b *Beneficiary) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}
