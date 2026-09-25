package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorPaymentMethodDoesNotExist = "error-payment-method-does-not-exist"

// ErrorPaymentMethodDoesNotExist is when user does not exist
type ErrorPaymentMethodDoesNotExist struct {
	ID string
}

// Error returns the error string
func (e *ErrorPaymentMethodDoesNotExist) Error() string {
	return errorPaymentMethodDoesNotExist
}

// ErrorType returns error type as string
func (e *ErrorPaymentMethodDoesNotExist) ErrorType() string {
	return errorPaymentMethodDoesNotExist
}

// Data returns data of the error
func (e *ErrorPaymentMethodDoesNotExist) Data() string {
	return fmt.Sprintf("%v", e.ID)
}

// Message returns tring message of error
func (e *ErrorPaymentMethodDoesNotExist) Message() string {
	return fmt.Sprintf("PaymentMethod: %v does not exist", e.ID)
}

// JSONError returns json of the error
func (e *ErrorPaymentMethodDoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorPaymentMethodDoesNotExist) HTTPCode() int {
	return http.StatusNotFound
}
