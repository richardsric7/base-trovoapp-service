package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidPaymentDestinationPublicKey = "error-invalid-payment-destination-public-key"

//ErrorInvalidPaymentDestinationPublicKey gets thrown when payment destination public key is invalid
type ErrorInvalidPaymentDestinationPublicKey struct {
}

//Error returns the error string
func (e *ErrorInvalidPaymentDestinationPublicKey) Error() string {
	return errorInvalidPaymentDestinationPublicKey
}

//ErrorType returns error type as string
func (e *ErrorInvalidPaymentDestinationPublicKey) ErrorType() string {
	return errorInvalidPaymentDestinationPublicKey
}

//Data returns data of the error
func (e *ErrorInvalidPaymentDestinationPublicKey) Data() string {
	return "destination"
}

//Message returns tring message of error
func (e *ErrorInvalidPaymentDestinationPublicKey) Message() string {
	return "Destination public key/address is invalid"
}

//JSONError returns json of the error
func (e *ErrorInvalidPaymentDestinationPublicKey) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorInvalidPaymentDestinationPublicKey) HTTPCode() int {
	return http.StatusBadRequest
}
