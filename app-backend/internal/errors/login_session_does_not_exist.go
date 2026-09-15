package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorLoginSessionDoesNotExist = "error-login-session-does-not-exist"

//ErrorLoginSessionDoesNotExist is when user does not exist
type ErrorLoginSessionDoesNotExist struct {
	Username string
}

//Error returns the error string
func (e *ErrorLoginSessionDoesNotExist) Error() string {
	return errorLoginSessionDoesNotExist
}

//ErrorType returns error type as string
func (e *ErrorLoginSessionDoesNotExist) ErrorType() string {
	return errorLoginSessionDoesNotExist
}

//Data returns data of the error
func (e *ErrorLoginSessionDoesNotExist) Data() string {
	return fmt.Sprintf("%v", e.Username)
}

//Message returns tring message of error
func (e *ErrorLoginSessionDoesNotExist) Message() string {
	return fmt.Sprintf("username: %v does not have pending login session", e.Username)
}

//JSONError returns json of the error
func (e *ErrorLoginSessionDoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorLoginSessionDoesNotExist) HTTPCode() int {
	return http.StatusNotFound
}
