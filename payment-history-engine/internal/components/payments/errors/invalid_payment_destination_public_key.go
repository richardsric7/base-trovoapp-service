package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidPaymentDestinationAddress = "error-invalid-payment-destination-public-key"

// ErrorInvalidPaymentDestinationAddress gets thrown when payment destination public key is invalid
type ErrorInvalidPaymentDestinationAddress struct {
}

// Error returns the error string
func (e *ErrorInvalidPaymentDestinationAddress) Error() string {
	return errorInvalidPaymentDestinationAddress
}

// ErrorType returns error type as string
func (e *ErrorInvalidPaymentDestinationAddress) ErrorType() string {
	return errorInvalidPaymentDestinationAddress
}

// Data returns data of the error
func (e *ErrorInvalidPaymentDestinationAddress) Data() string {
	return "destination"
}

// Message returns tring message of error
func (e *ErrorInvalidPaymentDestinationAddress) Message() string {
	return "Destination public key/address is invalid"
}

// JSONError returns json of the error
func (e *ErrorInvalidPaymentDestinationAddress) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorInvalidPaymentDestinationAddress) HTTPCode() int {
	return http.StatusBadRequest
}
