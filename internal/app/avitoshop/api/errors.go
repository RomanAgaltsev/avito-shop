package api

import (
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrBadRequest          = &ErrorResponse{StatusCode: 400, Message: "Bad request"}
	ErrEmptyItem           = &ErrorResponse{StatusCode: 400, Message: "Empty item"}
	ErrNotEnoughCoins      = &ErrorResponse{StatusCode: 400, Message: "Not enough coins"}
	ErrWrongLoginPassword  = &ErrorResponse{StatusCode: 401, Message: "Wrong login/password"}
	ErrNotFound            = &ErrorResponse{StatusCode: 404, Message: "Resource not found"}
	ErrMethodNotAllowed    = &ErrorResponse{StatusCode: 405, Message: "Method not allowed"}
	ErrLoginIsAlreadyTaken = &ErrorResponse{StatusCode: 409, Message: "Login has already been taken"}
)

type ErrorResponse struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status_text,omitempty"`
	Message    string `json:"errors"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}
func ErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 400,
		StatusText: "Bad request",
		Message:    err.Error(),
	}
}
func ServerErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 500,
		StatusText: "Internal server error",
		Message:    err.Error(),
	}
}
