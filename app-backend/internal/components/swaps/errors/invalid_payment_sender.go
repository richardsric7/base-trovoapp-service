package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidPaymentSender = "error-invalid-payment-sender"

// ErrorInvalidPaymentSender gets thrown when authenticated public key does not match supplied username in url
type ErrorInvalidPaymentSender struct {
}

// Error returns the error string
func (e *ErrorInvalidPaymentSender) Error() string {
	return errorInvalidPaymentSender
}

// ErrorType returns error type as string
func (e *ErrorInvalidPaymentSender) ErrorType() string {
	return errorInvalidPaymentSender
}

// Data returns data of the error
func (e *ErrorInvalidPaymentSender) Data() string {
	return "sender"
}

// Message returns tring message of error
func (e *ErrorInvalidPaymentSender) Message() string {
	return "Invalid username or email supplied. Neither the username nor email matches public key with header X-BANTUPAY-PUBLIC-KEY"
}

// JSONError returns json of the error
func (e *ErrorInvalidPaymentSender) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorInvalidPaymentSender) HTTPCode() int {
	return http.StatusBadRequest
}
