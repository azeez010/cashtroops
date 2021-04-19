package types

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type PaymentStatus string

const (
	INITIALIZED PaymentStatus = "INITIALIZED"
	DONE        PaymentStatus = "DONE"
)

type Payment struct {
	ID            uuid.UUID     `json:"id" gorm:"primary_key"`
	UserId        string        `json:"user_id"`
	Amount        int64         `json:"amount"` // In USD
	Coin          string        `json:"coin"`
	CoinPrice     int64         `json:"coin_price"` // In USD
	BeneficiaryId string        `json:"beneficiary_id"`
	AddressUsed   string        `json:"address_used"`
	Status        PaymentStatus `json:"status"`
	Ts            time.Time     `json:"ts"`
	KoboAmount    int64         `json:"kobo_amount"`
	UsdAmount     float64       `json:"usd_amount"`
	BtcAmount     float64       `json:"btc_amount"`
	TimeUpdated   time.Time     `json:"time_updated"`
}

type Balance struct {
	ID     uuid.UUID `json:"id" gorm:"primary_key"`
	UserId string    `json:"user_id"`
	Coin   string    `json:"coin"`
	Value  int64     `json:"value"`
	Ts     time.Time `json:"ts"`
}

type Address struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key"`
	UserId   string    `json:"user_id"`
	Public   string    `json:"public"`
	Private  string    `json:"private"`
	Provider string    `json:"provider"`
	Coin     string    `json:"coin"`
	Ts       time.Time `json:"ts"`
}

type Rate struct {
	ID           uuid.UUID `json:"id" gorm:"primary_key"`
	CurrencyPair string    `json:"currency_pair"`
	Value        int64     `json:"value"`
	Ts           time.Time `json:"ts"`
}

type Hook struct {
	ID      uuid.UUID `json:"id" gorm:"primary_key"`
	URL     string    `json:"url"`
	HookId  string    `json:"hook_id"`
	Address string    `json:"address"`
	Ts      time.Time `json:"ts"`
}

func (a *Address) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

func (p *Payment) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

func (b *Balance) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

func (b *Hook) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

func (r *Rate) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

func NewBalance(userId, coin string, value int64) *Balance {
	return &Balance{
		UserId: userId,
		Coin:   coin,
		Value:  value,
		Ts:     time.Now(),
	}
}

func NewPayment() *Payment {
	return &Payment{}
}
