package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidWallet = "error-invalid-wallet"

// ErrorInvalidWallet for invalid wallet
type ErrorInvalidWallet struct {
	PublicKey string
}

// Error returns the error string
func (e *ErrorInvalidWallet) Error() string {
	return errorInvalidWallet
}

// ErrorType returns error type as string
func (e *ErrorInvalidWallet) ErrorType() string {
	return errorInvalidWallet
}

// Data returns data of the error
func (e *ErrorInvalidWallet) Data() string {
	return errorInvalidWallet
}

// Message returns tring message of error
func (e *ErrorInvalidWallet) Message() string {
	return fmt.Sprintf("wallet address %v is invalid", e.PublicKey)
}

// JSONError returns json of the error
func (e *ErrorInvalidWallet) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": "address", "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorInvalidWallet) HTTPCode() int {
	return http.StatusBadRequest
}
