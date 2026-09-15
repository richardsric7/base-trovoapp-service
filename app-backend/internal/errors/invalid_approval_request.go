package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidRequest = "error-invalid-request"

// ErrorInvalidRequest for invalid wallet
type ErrorInvalidRequest struct {
	ID string
}

// Error returns the error string
func (e *ErrorInvalidRequest) Error() string {
	return errorInvalidRequest
}

// ErrorType returns error type as string
func (e *ErrorInvalidRequest) ErrorType() string {
	return errorInvalidRequest
}

// Data returns data of the error
func (e *ErrorInvalidRequest) Data() string {
	return errorInvalidRequest
}

// Message returns tring message of error
func (e *ErrorInvalidRequest) Message() string {
	return fmt.Sprintf("approval id %v is invalid", e.ID)
}

// JSONError returns json of the error
func (e *ErrorInvalidRequest) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": "address", "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorInvalidRequest) HTTPCode() int {
	return http.StatusBadRequest
}
