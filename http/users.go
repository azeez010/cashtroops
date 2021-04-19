package http

import (
	"encoding/json"
	"github.com/adigunhammedolalekan/cashtroops/ops"
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	userOps ops.UserOps
	logger  *logrus.Logger
}

func NewUserHandler(userOps ops.UserOps, logger *logrus.Logger) *UserHandler {
	return &UserHandler{userOps: userOps, logger: logger}
}

func (handler *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body := &types.CreateUserOpts{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		BadRequestResponse(w, r, "malformed request body")
		return
	}
	newUser, err := handler.userOps.CreateUser(body)
	if err != nil {
		handler.logger.WithError(err).Error("/user/new failed")
		Respond(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "user.created", Data: newUser})
}

func (handler *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequestResponse(w, r, "malformed request body")
		return
	}
	user, err := handler.userOps.AuthenticateUser(body.Email, body.Password)
	if err != nil {
		handler.logger.WithError(err).Error("failed to authenticate user")
		Respond(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "user.authenticated", Data: user})
}

func (handler *UserHandler) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	_, err := handler.userOps.GetSession(r.Header.Get(accountHeaderKey))
	if err != nil {
		ForbiddenRequestResponse(w, r, err.Error())
		return
	}
	var body struct {
		Code  string `json:"code"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequestResponse(w, r, "malformed request body")
		return
	}
	err = handler.userOps.ActivateAccount(body.Code, body.Email)
	if err != nil {
		handler.logger.WithError(err).Error("failed to activate account")
		Respond(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "account activated"})
}

func (handler *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	sess, err := handler.userOps.GetSession(r.Header.Get(accountHeaderKey))
	if err != nil {
		ForbiddenRequestResponse(w, r, err.Error())
		return
	}
	user, err := handler.userOps.GetUserByEmail(sess.Email)
	if err != nil {
		handler.logger.WithError(err).Error("user account not found")
		NotFoundResponse(w, r, "account does not exists")
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "success", Data: user})
}

func (handler *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	err := handler.userOps.RequestPasswordReset(email)
	if err != nil {
		handler.logger.WithError(err).Error("failed to request password reset")
		Respond(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "password reset requested"})
}

func (handler *UserHandler) VerifyPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code  string `json:"code"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequestResponse(w, r, "malformed request body")
		return
	}
	tk, err := handler.userOps.VerifyPasswordResetRequest(body.Code, body.Email)
	if err != nil {
		handler.logger.WithError(err).Error("cannot verify password reset details")
		Respond(w, r, err)
		return
	}
	type tempResponse struct {
		TokenId string `json:"token_id"`
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "code verified", Data: &tempResponse{TokenId: tk.ID.String()}})
}

func (handler *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NewPassword string `json:"new_password"`
		TokenId     string `json:"token_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequestResponse(w, r, "malformed request body")
		return
	}
	err := handler.userOps.ResetPassword(body.TokenId, body.NewPassword)
	if err != nil {
		handler.logger.WithError(err).Error("cannot verify password reset details")
		Respond(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "password changed"})
}

func (handler *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	sess, err := handler.userOps.GetSession(r.Header.Get(accountHeaderKey))
	if err != nil {
		ForbiddenRequestResponse(w, r, err.Error())
		return
	}
	var body struct {
		NewPassword string `json:"new_password"`
		OldPassword string `json:"old_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequestResponse(w, r, "malformed request body")
		return
	}

	err = handler.userOps.ChangePassword(sess.ID.String(), body.OldPassword, body.NewPassword)
	if err != nil {
		handler.logger.WithError(err).Error("user account not found")
		Respond(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "password changed"})
}
