package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorDestinationAddressCannotReceiveAsset = "error-destination-public-key-cannot-receive-asset"

// ErrorDestinationAddressCannotReceiveAsset gets thrown when payment destination public key does not trust asset
type ErrorDestinationAddressCannotReceiveAsset struct {
}

// Error returns the error string
func (e *ErrorDestinationAddressCannotReceiveAsset) Error() string {
	return errorDestinationAddressCannotReceiveAsset
}

// ErrorType returns error type as string
func (e *ErrorDestinationAddressCannotReceiveAsset) ErrorType() string {
	return errorDestinationAddressCannotReceiveAsset
}

// Data returns data of the error
func (e *ErrorDestinationAddressCannotReceiveAsset) Data() string {
	return "destination"
}

// Message returns tring message of error
func (e *ErrorDestinationAddressCannotReceiveAsset) Message() string {
	return "destination public key/address cannot receive this asset"
}

// JSONError returns json of the error
func (e *ErrorDestinationAddressCannotReceiveAsset) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorDestinationAddressCannotReceiveAsset) HTTPCode() int {
	return http.StatusBadRequest
}
