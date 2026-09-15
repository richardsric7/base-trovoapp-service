package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorDestinationPublicKeyCannotReceiveAsset = "error-destination-public-key-cannot-receive-asset"

//ErrorDestinationPublicKeyCannotReceiveAsset gets thrown when payment destination public key does not trust asset
type ErrorDestinationPublicKeyCannotReceiveAsset struct {
}

//Error returns the error string
func (e *ErrorDestinationPublicKeyCannotReceiveAsset) Error() string {
	return errorDestinationPublicKeyCannotReceiveAsset
}

//ErrorType returns error type as string
func (e *ErrorDestinationPublicKeyCannotReceiveAsset) ErrorType() string {
	return errorDestinationPublicKeyCannotReceiveAsset
}

//Data returns data of the error
func (e *ErrorDestinationPublicKeyCannotReceiveAsset) Data() string {
	return "destination"
}

//Message returns tring message of error
func (e *ErrorDestinationPublicKeyCannotReceiveAsset) Message() string {
	return "destination public key/address cannot receive this asset"
}

//JSONError returns json of the error
func (e *ErrorDestinationPublicKeyCannotReceiveAsset) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorDestinationPublicKeyCannotReceiveAsset) HTTPCode() int {
	return http.StatusBadRequest
}
