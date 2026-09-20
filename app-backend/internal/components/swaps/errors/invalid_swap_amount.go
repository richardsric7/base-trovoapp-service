package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidSwapAmount = "error-invalid-swap-amount"

// ErrorInvalidSwapAmount for an invalid payment amount
type ErrorInvalidSwapAmount struct {
}

// Error returns the error string
func (e *ErrorInvalidSwapAmount) Error() string {
	return errorInvalidSwapAmount
}

// ErrorType returns error type as string
func (e *ErrorInvalidSwapAmount) ErrorType() string {
	return errorInvalidSwapAmount
}

// Data returns data of the error
func (e *ErrorInvalidSwapAmount) Data() string {
	return "sourceAmount"
}

// Message returns tring message of error
func (e *ErrorInvalidSwapAmount) Message() string {
	return "Invalid swap amount"
}

// JSONError returns json of the error
func (e *ErrorInvalidSwapAmount) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorInvalidSwapAmount) HTTPCode() int {
	return http.StatusBadRequest
}
