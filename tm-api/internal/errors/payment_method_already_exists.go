package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorPaymentMethodAlreadyExists = "error-payment-method-already-exists"

// ErrorPaymentMethodAlreadyExists is returned during registration when an email is already in use
type ErrorPaymentMethodAlreadyExists struct {
	Detail string
}

// Error returns the error string
func (e *ErrorPaymentMethodAlreadyExists) Error() string {
	return errorPaymentMethodAlreadyExists
}

// ErrorType returns error type as string
func (e *ErrorPaymentMethodAlreadyExists) ErrorType() string {
	return errorPaymentMethodAlreadyExists
}

// Data returns data of the error
func (e *ErrorPaymentMethodAlreadyExists) Data() string {
	return "paymentMethodId"
}

// Message returns tring message of error
func (e *ErrorPaymentMethodAlreadyExists) Message() string {
	return e.Detail
}

// JSONError returns json of the error
func (e *ErrorPaymentMethodAlreadyExists) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorPaymentMethodAlreadyExists) HTTPCode() int {
	return http.StatusConflict
}
