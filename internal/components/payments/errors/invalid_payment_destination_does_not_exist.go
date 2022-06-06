package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorPaymentDestinationDoesNotExist = "error-payment-destination-does-not-exist"

//ErrorPaymentDestinationDoesNotExist gets thrown when payment destination username or email does not exist
type ErrorPaymentDestinationDoesNotExist struct {
}

//Error returns the error string
func (e *ErrorPaymentDestinationDoesNotExist) Error() string {
	return errorPaymentDestinationDoesNotExist
}

//ErrorType returns error type as string
func (e *ErrorPaymentDestinationDoesNotExist) ErrorType() string {
	return errorPaymentDestinationDoesNotExist
}

//Data returns data of the error
func (e *ErrorPaymentDestinationDoesNotExist) Data() string {
	return "destination"
}

//Message returns tring message of error
func (e *ErrorPaymentDestinationDoesNotExist) Message() string {
	return "Destination email or username does not exist"
}

//JSONError returns json of the error
func (e *ErrorPaymentDestinationDoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorPaymentDestinationDoesNotExist) HTTPCode() int {
	return http.StatusBadRequest
}
