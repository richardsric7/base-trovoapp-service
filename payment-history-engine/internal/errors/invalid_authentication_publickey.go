package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// errorInvalidAuthenticationAddress holds error of public key
const errorInvalidAuthenticationAddress = "error-invalid-authentication-public-key"

// ErrorInvalidAuthenticationAddress is struct for returning error
type ErrorInvalidAuthenticationAddress struct {
	Address string
}

// Error returns the error string
func (e *ErrorInvalidAuthenticationAddress) Error() string {
	return errorInvalidAddress
}

// ErrorType returns error type as string
func (e *ErrorInvalidAuthenticationAddress) ErrorType() string {
	return errorInvalidAuthenticationAddress
}

// Data returns data of the error
func (e *ErrorInvalidAuthenticationAddress) Data() string {
	return "X-TW-PUBLIC-KEY"
}

// Message returns tring message of error
func (e *ErrorInvalidAuthenticationAddress) Message() string {
	return "Invalid Header: [X-TW-PUBLIC-KEY= " + e.Address + "]"
}

// JSONError returns json of the error
func (e *ErrorInvalidAuthenticationAddress) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorInvalidAuthenticationAddress) HTTPCode() int {
	return http.StatusUnauthorized
}
