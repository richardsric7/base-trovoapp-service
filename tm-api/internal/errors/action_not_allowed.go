package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorActionNotAllowed = "error-action-not-allowed"

// ErrorActionNotAllowed is when Order does not exist
type ErrorActionNotAllowed struct {
	Param      string
	ErrMessage string
	Code       int
}

// Error returns the error string
func (e *ErrorActionNotAllowed) Error() string {
	return errorActionNotAllowed
}

// ErrorType returns error type as string
func (e *ErrorActionNotAllowed) ErrorType() string {
	return errorActionNotAllowed
}

// Data returns data of the error
func (e *ErrorActionNotAllowed) Data() string {
	return e.Param
}

// Message returns tring message of error
func (e *ErrorActionNotAllowed) Message() string {
	if len(e.ErrMessage) == 0 {
		return "action you are trying to perform is not allowed."
	}
	return e.ErrMessage
}

// JSONError returns json of the error
func (e *ErrorActionNotAllowed) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Param, "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorActionNotAllowed) HTTPCode() int {
	if e.Code >= 300 {
		return e.Code
	}
	return http.StatusBadRequest
}
