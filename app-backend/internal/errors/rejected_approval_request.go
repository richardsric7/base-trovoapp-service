package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorRejectedRequest = "error-rejected-request"

// ErrorRejectedRequest for invalid wallet
type ErrorRejectedRequest struct {
	ID string
}

// Error returns the error string
func (e *ErrorRejectedRequest) Error() string {
	return errorRejectedRequest
}

// ErrorType returns error type as string
func (e *ErrorRejectedRequest) ErrorType() string {
	return errorRejectedRequest
}

// Data returns data of the error
func (e *ErrorRejectedRequest) Data() string {
	return errorRejectedRequest
}

// Message returns tring message of error
func (e *ErrorRejectedRequest) Message() string {
	return "Approval request is already rejected. Only pending requests can be rejected."
}

// JSONError returns json of the error
func (e *ErrorRejectedRequest) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.ID, "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorRejectedRequest) HTTPCode() int {
	return http.StatusBadRequest
}
