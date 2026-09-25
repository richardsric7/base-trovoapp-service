package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorOperationFailed = "error-operation-failed"

// ErrorOperationFailed for invalid email formats
type ErrorOperationFailed struct {
	Param      string
	ErrMessage string
	Code       int
}

// Error returns the error string
func (e *ErrorOperationFailed) Error() string {
	return errorOperationFailed
}

// ErrorType returns error type as string
func (e *ErrorOperationFailed) ErrorType() string {
	return errorOperationFailed
}

// Data returns data of the error
func (e *ErrorOperationFailed) Data() string {
	return e.Param
}

// Message returns tring message of error
func (e *ErrorOperationFailed) Message() string {
	return e.ErrMessage
}

// JSONError returns json of the error
func (e *ErrorOperationFailed) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Param, "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorOperationFailed) HTTPCode() int {
	if e.Code >= 300 {
		return e.Code
	}
	return http.StatusBadRequest

}
