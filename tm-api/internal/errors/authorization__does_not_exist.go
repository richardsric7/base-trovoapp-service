package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorAuthorizationoesNotExist = "error-authorization-does-not-exist"

// ErrorAuthorizationoesNotExist is when user does not exist
type ErrorAuthorizationoesNotExist struct {
	Username string
}

// Error returns the error string
func (e *ErrorAuthorizationoesNotExist) Error() string {
	return errorAuthorizationoesNotExist
}

// ErrorType returns error type as string
func (e *ErrorAuthorizationoesNotExist) ErrorType() string {
	return errorAuthorizationoesNotExist
}

// Data returns data of the error
func (e *ErrorAuthorizationoesNotExist) Data() string {
	return fmt.Sprintf("%v", e.Username)
}

// Message returns tring message of error
func (e *ErrorAuthorizationoesNotExist) Message() string {
	return fmt.Sprintf("username: %v does not have pending login session", e.Username)
}

// JSONError returns json of the error
func (e *ErrorAuthorizationoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorAuthorizationoesNotExist) HTTPCode() int {
	return http.StatusNotFound
}
