package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidAddress = "error-invalid-public-key"

// ErrorInvalidAddress is error struct
type ErrorInvalidAddress struct {
	Address string
}

// Error returns the error string
func (e *ErrorInvalidAddress) Error() string {
	return errorInvalidAddress
}

// ErrorType returns error type as string
func (e *ErrorInvalidAddress) ErrorType() string {
	return errorInvalidAddress
}

// Data returns data of the error
func (e *ErrorInvalidAddress) Data() string {
	return errorInvalidAddress
}

// Message returns tring message of error
func (e *ErrorInvalidAddress) Message() string {
	return "Invalid Public Key [" + e.Address + "]"
}

// JSONError returns json of the error
func (e *ErrorInvalidAddress) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": "publicKey", "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorInvalidAddress) HTTPCode() int {
	return http.StatusBadRequest
}
