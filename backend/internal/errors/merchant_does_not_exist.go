package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorServiceDoesNotExist = "error-service-does-not-exist"

// ErrorServiceDoesNotExist is when user does not exist
type ErrorServiceDoesNotExist struct {
	Username string
}

// Error returns the error string
func (e *ErrorServiceDoesNotExist) Error() string {
	return errorServiceDoesNotExist
}

// ErrorType returns error type as string
func (e *ErrorServiceDoesNotExist) ErrorType() string {
	return errorServiceDoesNotExist
}

// Data returns data of the error
func (e *ErrorServiceDoesNotExist) Data() string {
	return fmt.Sprintf("%v", e.Username)
}

// Message returns tring message of error
func (e *ErrorServiceDoesNotExist) Message() string {
	return fmt.Sprintf("service: %v does not exist", e.Username)
}

// JSONError returns json of the error
func (e *ErrorServiceDoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorServiceDoesNotExist) HTTPCode() int {
	return http.StatusNotFound
}
