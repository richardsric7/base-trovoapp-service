package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorOfferAlreadyExists = "error-offer-already-exists"

// ErrorOfferAlreadyExists is returned during registration when an email is already in use
type ErrorOfferAlreadyExists struct {
	Detail string
}

// Error returns the error string
func (e *ErrorOfferAlreadyExists) Error() string {
	return errorOfferAlreadyExists
}

// ErrorType returns error type as string
func (e *ErrorOfferAlreadyExists) ErrorType() string {
	return errorOfferAlreadyExists
}

// Data returns data of the error
func (e *ErrorOfferAlreadyExists) Data() string {
	return "OfferId"
}

// Message returns tring message of error
func (e *ErrorOfferAlreadyExists) Message() string {
	return e.Detail
}

// JSONError returns json of the error
func (e *ErrorOfferAlreadyExists) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorOfferAlreadyExists) HTTPCode() int {
	return http.StatusConflict
}
