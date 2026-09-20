package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorSourceAndDestinationAssetAreSame = "error-source-and-destination-asset-are-same"

//ErrorSourceAndDestinationAssetAreSame for an invalid payment amount
type ErrorSourceAndDestinationAssetAreSame struct {
}

//Error returns the error string
func (e *ErrorSourceAndDestinationAssetAreSame) Error() string {
	return errorSourceAndDestinationAssetAreSame
}

//ErrorType returns error type as string
func (e *ErrorSourceAndDestinationAssetAreSame) ErrorType() string {
	return errorSourceAndDestinationAssetAreSame
}

//Data returns data of the error
func (e *ErrorSourceAndDestinationAssetAreSame) Data() string {
	return "destinationAsset"
}

//Message returns tring message of error
func (e *ErrorSourceAndDestinationAssetAreSame) Message() string {
	return "source and destination cannot be same"
}

//JSONError returns json of the error
func (e *ErrorSourceAndDestinationAssetAreSame) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorSourceAndDestinationAssetAreSame) HTTPCode() int {
	return http.StatusBadRequest
}
