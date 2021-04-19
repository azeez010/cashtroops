package ops

import (
	"github.com/adigunhammedolalekan/cashtroops/errors"
	"github.com/adigunhammedolalekan/cashtroops/fn"
	"github.com/adigunhammedolalekan/cashtroops/session"
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	passwordResetTokenTable = "password_reset_tokens"
)

type UserOps interface {
	CreateUser(user *types.CreateUserOpts) (*types.User, error)
	AuthenticateUser(email, password string) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	ActivateAccount(code, email string) error
	GetSession(key string) (*types.User, error)
	RequestPasswordReset(email string) error
	VerifyPasswordResetRequest(code, email string) (*types.PasswordResetToken, error)
	ResetPassword(tokenId, newPassword string) error
	GetPasswordResetToken(code, email string) (*types.PasswordResetToken, error)
	GetPasswordResetTokenById(id string) (*types.PasswordResetToken, error)
	ChangePassword(userId, oldPassword, newPassword string) error
	GetUserByAttr(attr string, value interface{}) (*types.User, error)
}

type userOps struct {
	db      *gorm.DB
	session session.Store
	logger  *logrus.Logger
}

func NewUserOps(db *gorm.DB, sess session.Store, logger *logrus.Logger) UserOps {
	return &userOps{
		db:      db,
		session: sess,
		logger:  logger,
	}
}

func (u *userOps) CreateUser(user *types.CreateUserOpts) (*types.User, error) {
	if err := fn.ValidateEmail(user.Email); err != nil {
		return nil, errors.New(http.StatusBadRequest, err.Error())
	}
	existing, err := u.GetUserByEmail(user.Email)
	if err == nil && existing.Email != "" {
		return nil, errors.New(http.StatusConflict, "email is already in use by another customer")
	}

	if user.FirstName == "" && user.LastName == "" {
		return nil, errors.New(http.StatusBadRequest, "first name or last name must be present")
	}
	if err := fn.ValidatePassword(user.Password); err != nil {
		return nil, errors.New(http.StatusBadRequest, err.Error())
	}
	user.Password = fn.HashPassword(user.Password)
	newUser := types.NewUser(user)
	token := types.NewToken(fn.GenerateRandomString(64), newUser)
	newUser.Token = token.Key

	if err := u.session.Create(token); err != nil {
		u.logger.WithError(err).Error("failed to create auth token for new user")
		return nil, errors.New(http.StatusInternalServerError, "failed to create account at this time. please retry later")
	}
	tx := u.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	if err := tx.Table("users").Create(newUser).Error; err != nil {
		tx.Rollback()
		u.logger.WithError(err).Error("failed to create user in the database")
		return nil, errors.New(http.StatusInternalServerError, "failed to create account at this time. please retry later")
	}
	verification := types.NewVerification(user.Email, fn.GenRandomCode())
	if err := tx.Table("verifications").Create(verification).Error; err != nil {
		tx.Rollback()
		u.logger.WithError(err).Error("failed to create verification")
		return nil, errors.New(http.StatusInternalServerError, "failed to create account at this time. please retry later")
	}
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(http.StatusInternalServerError, "failed to create account at this time. please retry later")
	}

	go func(code, email, user string, logger *logrus.Logger) {
		mailBody, err := fn.GenerateWelcomeEmail(user, code)
		if err != nil {
			logger.WithError(err).Error("failed to generate welcome email")
			return
		}
		if err := fn.SendEmail(&types.MailRequest{
			User:  user,
			Email: email,
			Title: "Welcome to CashTroops",
			Body:  mailBody,
		}); err != nil {
			logger.WithError(err).Error("failed to send email")
		}
	}(verification.Code, newUser.Email, newUser.Name(), u.logger)
	return newUser, nil
}

func (u *userOps) AuthenticateUser(email, password string) (*types.User, error) {
	if err := fn.ValidateEmail(email); err != nil {
		return nil, errors.New(http.StatusBadRequest, err.Error())
	}

	if err := fn.ValidatePassword(password); err != nil {
		return nil, errors.New(http.StatusBadRequest, err.Error())
	}
	user, err := u.GetUserByEmail(email)
	if err != nil {
		u.logger.WithField("email", email).WithError(err).Error("user not found")
		return nil, errors.New(http.StatusForbidden, "email and password combination does not match")
	}
	if ok := fn.VerifyHashPassword(user.Password, password); !ok {
		return nil, errors.New(http.StatusForbidden, "email and password combination does not match")
	}
	token := types.NewToken(fn.GenerateRandomString(64), user)
	user.Token = token.Key
	if err := u.session.Create(token); err != nil {
		u.logger.WithError(err).Error("failed to create auth token for new user")
		return nil, errors.New(http.StatusInternalServerError, "failed to create account at this time. please retry later")
	}
	return user, nil
}

func (u *userOps) GetUserByEmail(email string) (*types.User, error) {
	return u.GetUserByAttr("email", email)
}

func (u *userOps) ActivateAccount(code, email string) error {
	v := &types.Verification{}
	err := u.db.Table("verifications").Where("code = ? AND email = ?", code, email).First(v).Error
	if err != nil {
		return err
	}
	return u.db.Table("verifications").Where("code = ? AND email = ?", code, email).UpdateColumn("activated", true).Error
}

func (u *userOps) GetSession(key string) (*types.User, error) {
	return u.session.Get(key)
}

func (u *userOps) RequestPasswordReset(email string) error {
	user, err := u.GetUserByEmail(email)
	if err != nil {
		return err
	}
	tk := types.NewPasswordResetToken(fn.GenRandomCode(), email, user.ID.String())
	if err := u.db.Table(passwordResetTokenTable).Create(tk).Error; err != nil {
		u.logger.WithError(err).Error("failed to create password reset token")
		return errors.New(http.StatusInternalServerError, "failed to reset password at this time. please retry")
	}

	go func(code, email string) {
		emailBody, err := fn.GenerateResetPasswordEmail(code)
		if err != nil {
			u.logger.WithError(err).Error("failed to generate password reset email")
			return
		}
		if err := fn.SendEmail(&types.MailRequest{
			User:  "",
			Email: email,
			Title: "Password Reset Instructions - CashTroops",
			Body:  emailBody,
		}); err != nil {
			u.logger.WithError(err).Error("failed to send password reset email")
		}
	}(tk.Code, user.Email)
	return nil
}

func (u *userOps) VerifyPasswordResetRequest(code, email string) (*types.PasswordResetToken, error) {
	tk, err := u.GetPasswordResetToken(code, email)
	if err != nil {
		u.logger.WithError(err).Error("failed to get password reset token")
		return nil, errors.New(http.StatusUnauthorized, "password reset request not found")
	}
	return tk, u.db.Table(passwordResetTokenTable).Where("id = ?", tk.ID).UpdateColumn("verified", true).Error
}

func (u *userOps) ResetPassword(tokenId, newPassword string) error {
	tk, err := u.GetPasswordResetTokenById(tokenId)
	if err != nil {
		u.logger.WithError(err).Error("failed to get password reset token")
		return errors.New(http.StatusUnauthorized, "password reset request not found")
	}
	if err := fn.ValidatePassword(newPassword); err != nil {
		return errors.New(http.StatusBadRequest, err.Error())
	}
	newHashedPassword := fn.HashPassword(newPassword)
	return u.db.Table("users").Where("id = ?", tk.OwnerId).UpdateColumn("password", newHashedPassword).Error
}

func (u *userOps) GetPasswordResetToken(code, email string) (*types.PasswordResetToken, error) {
	tk := &types.PasswordResetToken{}
	err := u.db.Table(passwordResetTokenTable).Where("code = ? AND email = ?", code, email).First(tk).Error
	return tk, err
}

func (u *userOps) GetPasswordResetTokenById(id string) (*types.PasswordResetToken, error) {
	tk := &types.PasswordResetToken{}
	err := u.db.Table(passwordResetTokenTable).Where("id = ?", id).First(tk).Error
	return tk, err
}

func (u *userOps) ChangePassword(userId, old, newPassword string) error {
	user, err := u.GetUserByAttr("id", userId)
	if err != nil {
		u.logger.WithError(err).Error("failed to get userById")
		return errors.New(http.StatusNotFound, "user not found")
	}
	if ok := fn.VerifyHashPassword(user.Password, old); !ok {
		return errors.New(http.StatusUnauthorized, "old password does not match our record")
	}
	if err := fn.ValidatePassword(newPassword); err != nil {
		return errors.New(http.StatusBadRequest, err.Error())
	}
	return u.db.Table("users").Where("id = ?", user.ID).UpdateColumn("password", fn.HashPassword(newPassword)).Error
}

func (u *userOps) GetUserByAttr(attr string, value interface{}) (*types.User, error) {
	user := &types.User{}
	err := u.db.Table("users").Where(attr+" = ?", value).First(user).Error
	return user, err
}
