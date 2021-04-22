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

type AccountHandler struct {
	accountOps ops.AccountOps
	userOps    ops.UserOps
	logger     *logrus.Logger
}

func NewAccountHandler(accountOps ops.AccountOps, userOps ops.UserOps, logger *logrus.Logger) *AccountHandler {
	return &AccountHandler{accountOps: accountOps, userOps: userOps, logger: logger}
}

func (handler *AccountHandler) AddBeneficiary(w http.ResponseWriter, r *http.Request) {
	sess, err := handler.userOps.GetSession(r.Header.Get(accountHeaderKey))
	if err != nil {
		ForbiddenRequestResponse(w, r, err.Error())
		return
	}
	body := &types.CreateBeneficiaryOpts{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		BadRequestResponse(w, r, "malformed request body")
		return
	}
	bf, err := handler.accountOps.AddBeneficiary(sess.ID.String(), body)
	if err != nil {
		handler.logger.WithError(err).Error("/beneficiary/new failed")
		Respond(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "beneficiary added", Data: bf})
}

func (handler *AccountHandler) RemoveBeneficiary(w http.ResponseWriter, r *http.Request) {
	sess, err := handler.userOps.GetSession(r.Header.Get(accountHeaderKey))
	if err != nil {
		ForbiddenRequestResponse(w, r, err.Error())
		return
	}
	bfId := chi.URLParam(r, "id")
	err = handler.accountOps.RemoveBeneficiary(sess.ID.String(), bfId)
	if err != nil {
		handler.logger.WithError(err).Error("/beneficiary/id/remove failed")
		Respond(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "beneficiary removed"})
}

func (handler *AccountHandler) ListBeneficiaries(w http.ResponseWriter, r *http.Request) {
	sess, err := handler.userOps.GetSession(r.Header.Get(accountHeaderKey))
	if err != nil {
		ForbiddenRequestResponse(w, r, err.Error())
		return
	}

	data, err := handler.accountOps.ListBeneficiariesFor(sess.ID.String())
	if err != nil {
		handler.logger.WithError(err).Error("/me/beneficiaries failed")
		InternalServerErrorResponse(w, r, "failed to fetch beneficiaries at the moment. please retry")
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "success", Data: data})
}

func (handler *AccountHandler) Banks(w http.ResponseWriter, r *http.Request) {
	sess, err := handler.userOps.GetSession(r.Header.Get(accountHeaderKey))
	if err != nil {
		ForbiddenRequestResponse(w, r, err.Error())
		return
	}
	handler.logger.WithField("account_id", sess.ID.String()).
		Info("requesting bank list")
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "success", Data: handler.accountOps.GetBanks()})
}
