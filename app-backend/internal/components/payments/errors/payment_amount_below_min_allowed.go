package errors

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const errorPaymentAmountBelowMinAllowed = "error-payment-amount-below-min-allowed"

// ErrorPaymentAmountBelowMinAllowed for an invalid payment amount
type ErrorPaymentAmountBelowMinAllowed struct {
}

// Error returns the error string
func (e *ErrorPaymentAmountBelowMinAllowed) Error() string {
	return errorPaymentAmountBelowMinAllowed
}

// ErrorType returns error type as string
func (e *ErrorPaymentAmountBelowMinAllowed) ErrorType() string {
	return errorPaymentAmountBelowMinAllowed
}

// Data returns data of the error
func (e *ErrorPaymentAmountBelowMinAllowed) Data() string {
	return "amount"
}

// Message returns tring message of error
func (e *ErrorPaymentAmountBelowMinAllowed) Message() string {
	minAmountSendable := "0.0000001"
	if os.Getenv("MIN_SENDABLE_AMOUNT") != "" && os.Getenv("MIN_SENDABLE_AMOUNT") != "0" {
		minAmountSendable = os.Getenv("MIN_SENDABLE_AMOUNT")
	}
	return "amount is below minimum allowed: " + minAmountSendable
}

// JSONError returns json of the error
func (e *ErrorPaymentAmountBelowMinAllowed) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorPaymentAmountBelowMinAllowed) HTTPCode() int {
	return http.StatusBadRequest
}
