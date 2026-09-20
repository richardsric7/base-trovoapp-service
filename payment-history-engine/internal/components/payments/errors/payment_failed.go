package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorPaymentFailed = "error-payment-failed"

// ErrorPaymentFailed for when generated transaction does not match supplied transaction
type ErrorPaymentFailed struct {
}

// Error returns the error string
func (e *ErrorPaymentFailed) Error() string {
	return errorPaymentFailed
}

// ErrorType returns error type as string
func (e *ErrorPaymentFailed) ErrorType() string {
	return errorPaymentFailed
}

// Data returns data of the error
func (e *ErrorPaymentFailed) Data() string {
	return "transaction"
}

// Message returns tring message of error
func (e *ErrorPaymentFailed) Message() string {
	return "Transaction failed. Please try again later"
}

// JSONError returns json of the error
func (e *ErrorPaymentFailed) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorPaymentFailed) HTTPCode() int {
	return http.StatusInternalServerError
}
