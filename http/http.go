package http

import (
	"github.com/adigunhammedolalekan/cashtroops/errors"
	"github.com/go-chi/render"
	"net/http"
)

var (
	accountHeaderKey = "X-Account-Token"
)

type SuccessResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusBadRequest)
	render.Respond(w, r, &ErrorResponse{Error: true, Message: message})
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.Respond(w, r, &ErrorResponse{Error: true, Message: message})
}

func ForbiddenRequestResponse(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusForbidden)
	render.Respond(w, r, &ErrorResponse{Error: true, Message: message})
}

func InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusForbidden)
	render.Respond(w, r, &ErrorResponse{Error: true, Message: message})
}

func Respond(w http.ResponseWriter, r *http.Request, err error) {
	if e, ok := err.(*errors.Error); ok {
		render.Status(r, e.Code)
		render.Respond(w, r, &ErrorResponse{Error: true, Message: e.Message})
	} else {
		InternalServerErrorResponse(w, r, err.Error())
	}
}
