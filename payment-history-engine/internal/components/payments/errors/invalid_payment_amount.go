package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidPaymentAmount = "error-invalid-payment-amount"

//ErrorInvalidPaymentAmount for an invalid payment amount
type ErrorInvalidPaymentAmount struct {
}

//Error returns the error string
func (e *ErrorInvalidPaymentAmount) Error() string {
	return errorInvalidPaymentAmount
}

//ErrorType returns error type as string
func (e *ErrorInvalidPaymentAmount) ErrorType() string {
	return errorInvalidPaymentAmount
}

//Data returns data of the error
func (e *ErrorInvalidPaymentAmount) Data() string {
	return "amount"
}

//Message returns tring message of error
func (e *ErrorInvalidPaymentAmount) Message() string {
	return "Invalid payment amount"
}

//JSONError returns json of the error
func (e *ErrorInvalidPaymentAmount) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorInvalidPaymentAmount) HTTPCode() int {
	return http.StatusBadRequest
}
