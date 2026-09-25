package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorAssetDoesNotExist = "error-asset-does-not-exist"

// ErrorAssetDoesNotExist is when Offer does not exist
type ErrorAssetDoesNotExist struct {
	AssetID string
}

// Error returns the error string
func (e *ErrorAssetDoesNotExist) Error() string {
	return errorAssetDoesNotExist
}

// ErrorType returns error type as string
func (e *ErrorAssetDoesNotExist) ErrorType() string {
	return errorAssetDoesNotExist
}

// Data returns data of the error
func (e *ErrorAssetDoesNotExist) Data() string {
	return fmt.Sprintf("%v", e.AssetID)
}

// Message returns tring message of error
func (e *ErrorAssetDoesNotExist) Message() string {
	return fmt.Sprintf("Asset: %v does not exist", e.AssetID)
}

// JSONError returns json of the error
func (e *ErrorAssetDoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorAssetDoesNotExist) HTTPCode() int {
	return http.StatusNotFound
}
