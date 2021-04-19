package http

import (
	"encoding/json"
	"github.com/adigunhammedolalekan/cashtroops/libs/bc"
	"github.com/adigunhammedolalekan/cashtroops/ops"
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/blockcypher/gobcy"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PaymentHandler struct {
	paymentOps ops.PaymentOps
	userOps    ops.UserOps
	logger     *logrus.Logger
}

func NewPaymentHandler(paymentOps ops.PaymentOps, userOps ops.UserOps, logger *logrus.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentOps: paymentOps,
		userOps:    userOps,
		logger:     logger,
	}
}

func (handler *PaymentHandler) InitializePayment(w http.ResponseWriter, r *http.Request) {
	sess, err := handler.userOps.GetSession(r.Header.Get(accountHeaderKey))
	if err != nil {
		ForbiddenRequestResponse(w, r, err.Error())
		return
	}
	body := &types.InitPaymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		handler.logger.WithError(err).Error("malformed json body")
		BadRequestResponse(w, r, "malformed request body")
		return
	}
	if body.AmountInt() <= 0 {
		handler.logger.WithField("amount", body.AmountInt()).Error("invalid amount supplied")
		BadRequestResponse(w, r, "The amount seems to be invalid")
		return
	}
	resp, err := handler.paymentOps.InitializePayment(sess.ID.String(), body)
	if err != nil {
		Respond(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "payment initialized", Data: resp})
}

func (handler *PaymentHandler) TxnEventHandler(w http.ResponseWriter, r *http.Request) {
	event, webHookId := r.Header.Get("X-EventType"), r.Header.Get("X-EventId")
	handler.logger.WithFields(logrus.Fields{
		"event_type": event,
		"webhook_id": webHookId,
	}).Info("recv new webhook")
	if event != bc.EventTypeTxConfirmation {
		BadRequestResponse(w, r, "unwanted event")
		return
	}
	if err := handler.paymentOps.VerifyWebHookId(webHookId); err != nil {
		handler.logger.WithError(err).Error("failed to verify webhookId")
		Respond(w, r, err)
		return
	}
	tx := &gobcy.TX{}
	if err := json.NewDecoder(r.Body).Decode(tx); err != nil {
		BadRequestResponse(w, r, "malformed request body")
		return
	}
	if tx.Confirmations < 2 {
		handler.logger.WithField("tx_confirmation", tx.Confirmations).Info("need 2 confirmations to proceed")
		render.Status(r, http.StatusOK)
		return
	}
	recvAddress := ""
	var amount int64 = 0
	if len(tx.Outputs) > 0 && len(tx.Outputs[0].Addresses) > 0 {
		recvAddress = tx.Outputs[0].Addresses[0]
		amount = tx.Outputs[0].Value.Int64()
	}
	handler.logger.WithFields(logrus.Fields{
		"recv_address": recvAddress,
		"amount":       amount,
	}).Info("recv new BTC")
	if recvAddress == "" || amount == 0 {
		BadRequestResponse(w, r, "invalid transaction data")
		return
	}
	err := handler.paymentOps.FinalizePayment(recvAddress, amount)
	if err != nil {
		Respond(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
}

func (handler *PaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	sess, err := handler.userOps.GetSession(r.Header.Get(accountHeaderKey))
	if err != nil {
		ForbiddenRequestResponse(w, r, err.Error())
		return
	}
	data, err := handler.paymentOps.ListPayments(sess.ID.String())
	if err != nil {
		NotFoundResponse(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &SuccessResponse{Error: false, Message: "success", Data: data})
}
