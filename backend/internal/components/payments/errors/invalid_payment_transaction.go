package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidPaymentTransaction = "error-invalid-payment-transaction"

//ErrorInvalidPaymentAmount for an invalid payment amount
type ErrorInvalidPaymentTransaction struct {
}

//Error returns the error string
func (e *ErrorInvalidPaymentTransaction) Error() string {
	return errorInvalidPaymentTransaction
}

//ErrorType returns error type as string
func (e *ErrorInvalidPaymentTransaction) ErrorType() string {
	return errorInvalidPaymentTransaction
}

//Data returns data of the error
func (e *ErrorInvalidPaymentTransaction) Data() string {
	return "transaction"
}

//Message returns tring message of error
func (e *ErrorInvalidPaymentTransaction) Message() string {
	return "Invalid base64 payment transaction"
}

//JSONError returns json of the error
func (e *ErrorInvalidPaymentTransaction) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorInvalidPaymentTransaction) HTTPCode() int {
	return http.StatusBadRequest
}
