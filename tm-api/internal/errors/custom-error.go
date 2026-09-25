package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CustomError is when user does not exist
type CustomError struct {
	Param      string
	Err        string
	ErrMessage string
	Code       int
}

// Error returns the error string
func (e *CustomError) Error() string {
	return strings.ToLower(strings.ReplaceAll(e.Err, " ", "-"))
}

// ErrorType returns error type as string
func (e *CustomError) ErrorType() string {
	return e.Error()
}

// Data returns data of the error
func (e *CustomError) Data() string {
	return fmt.Sprintf("%v", e.Param)
}

// Message returns tring message of error
func (e *CustomError) Message() string {
	return e.ErrMessage
}

// JSONError returns json of the error
func (e *CustomError) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *CustomError) HTTPCode() int {
	if e.Code >= 300 {
		return e.Code
	}
	return http.StatusBadRequest
}
