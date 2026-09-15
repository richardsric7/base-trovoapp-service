package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorOnlyViewAccessAllowedInPrimaryWallet = "error-only-view-access-allowed-in-primary-wallet"

// ErrorOnlyViewAccessAllowedInPrimaryWallet
type ErrorOnlyViewAccessAllowedInPrimaryWallet struct {
	Detail string
}

// Error returns the error string
func (e *ErrorOnlyViewAccessAllowedInPrimaryWallet) Error() string {
	return errorOnlyViewAccessAllowedInPrimaryWallet
}

// ErrorType returns error type as string
func (e *ErrorOnlyViewAccessAllowedInPrimaryWallet) ErrorType() string {
	return errorOnlyViewAccessAllowedInPrimaryWallet
}

// Data returns data of the error
func (e *ErrorOnlyViewAccessAllowedInPrimaryWallet) Data() string {
	return "accessList"
}

// Message returns tring message of error
func (e *ErrorOnlyViewAccessAllowedInPrimaryWallet) Message() string {
	if len(e.Detail) > 0 {
		return e.Detail
	}
	return "Only view access allowed in primary wallet."

}

// JSONError returns json of the error
func (e *ErrorOnlyViewAccessAllowedInPrimaryWallet) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorOnlyViewAccessAllowedInPrimaryWallet) HTTPCode() int {
	return http.StatusConflict
}
