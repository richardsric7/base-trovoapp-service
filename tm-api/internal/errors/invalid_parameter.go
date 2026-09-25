package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidParameter = "error-invalid-parameter"

// ErrorInvalidEmailFormat for invalid email formats
type ErrorInvalidParameter struct {
	Param      string
	ErrMessage string
	Code       int
}

// Error returns the error string
func (e *ErrorInvalidParameter) Error() string {
	return errorInvalidParameter
}

// ErrorType returns error type as string
func (e *ErrorInvalidParameter) ErrorType() string {
	return errorInvalidParameter
}

// Data returns data of the error
func (e *ErrorInvalidParameter) Data() string {
	return e.Param
}

// Message returns tring message of error
func (e *ErrorInvalidParameter) Message() string {
	return e.ErrMessage
}

// JSONError returns json of the error
func (e *ErrorInvalidParameter) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Param, "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorInvalidParameter) HTTPCode() int {
	if e.Code >= 300 {
		return e.Code
	}
	return http.StatusBadRequest

}
