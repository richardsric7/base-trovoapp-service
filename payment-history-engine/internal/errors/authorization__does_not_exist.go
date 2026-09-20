package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorAuthorizationDoesNotExist = "error-authorization-does-not-exist"

// ErrorAuthorizationDoesNotExist is when user does not exist
type ErrorAuthorizationDoesNotExist struct {
	Username string
}

// Error returns the error string
func (e *ErrorAuthorizationDoesNotExist) Error() string {
	return errorAuthorizationDoesNotExist
}

// ErrorType returns error type as string
func (e *ErrorAuthorizationDoesNotExist) ErrorType() string {
	return errorAuthorizationDoesNotExist
}

// Data returns data of the error
func (e *ErrorAuthorizationDoesNotExist) Data() string {
	return fmt.Sprintf("%v", e.Username)
}

// Message returns tring message of error
func (e *ErrorAuthorizationDoesNotExist) Message() string {
	return fmt.Sprintf("username: %v does not have pending login session", e.Username)
}

// JSONError returns json of the error
func (e *ErrorAuthorizationDoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorAuthorizationDoesNotExist) HTTPCode() int {
	return http.StatusNotFound
}
