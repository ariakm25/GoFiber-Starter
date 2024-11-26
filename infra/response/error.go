package response

import (
	"net/http"
)

type Error struct {
	Message  string
	Code     string
	HttpCode int
}

func NewError(msg string, code string, httpCode int) Error {
	return Error{
		Message:  msg,
		Code:     code,
		HttpCode: httpCode,
	}
}

func (e Error) Error() string {
	return e.Message
}

// General Error
var (
	ErrorInternal            = NewError("internal server error", "500", http.StatusInternalServerError)
	ErrorBadRequest          = NewError("bad request", "400", http.StatusBadRequest)
	ErrorUnauthorized        = NewError("unauthorized", "401", http.StatusUnauthorized)
	ErrorForbiddenAccess     = NewError("forbidden", "403", http.StatusForbidden)
	ErrorNotFound            = NewError("not found", "404", http.StatusNotFound)
	ErrorUnprocessableEntity = NewError("unprocessable entity", "422", http.StatusUnprocessableEntity)
)
