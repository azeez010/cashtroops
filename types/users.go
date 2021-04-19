package types

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"-"`
	Token     string    `json:"token" gorm:"-" sql:"-"`
	Ts        time.Time `json:"ts"`
}

type Token struct {
	Key     string    `json:"key"`
	Created time.Time `json:"created"`
	User    *User     `json:"user"`
}

type Verification struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	Code      string    `json:"code"`
	Email     string    `json:"email"`
	Activated bool      `json:"activated"`
	Ts        time.Time `json:"ts"`
}

type PasswordResetToken struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key"`
	Code     string    `json:"code"`
	Email    string    `json:"email"`
	OwnerId  string    `json:"owner"`
	Verified bool      `json:"verified"`
	Ts       time.Time `json:"ts"`
}

func (ps *PasswordResetToken) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

func (v *Verification) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

func (user *User) Name() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func NewToken(key string, account *User) *Token {
	return &Token{Key: key, User: account, Created: time.Now()}
}

func NewUser(user *CreateUserOpts) *User {
	return &User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Ts:        time.Now(),
	}
}

func NewVerification(email, code string) *Verification {
	return &Verification{
		Code:      code,
		Email:     email,
		Activated: false,
		Ts:        time.Now(),
	}
}

func NewPasswordResetToken(code, email, owner string) *PasswordResetToken {
	return &PasswordResetToken{
		Code:     code,
		Email:    email,
		OwnerId:  owner,
		Verified: false,
		Ts:       time.Now(),
	}
}
