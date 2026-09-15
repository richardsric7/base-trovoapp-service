package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorCompletedRequest = "error-completed-request"

// ErrorCompletedRequest for invalid wallet
type ErrorCompletedRequest struct {
	ID string
}

// Error returns the error string
func (e *ErrorCompletedRequest) Error() string {
	return errorCompletedRequest
}

// ErrorType returns error type as string
func (e *ErrorCompletedRequest) ErrorType() string {
	return errorCompletedRequest
}

// Data returns data of the error
func (e *ErrorCompletedRequest) Data() string {
	return errorCompletedRequest
}

// Message returns tring message of error
func (e *ErrorCompletedRequest) Message() string {
	return "Approval request is already completed. Only pending requests can be approved."
}

// JSONError returns json of the error
func (e *ErrorCompletedRequest) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.ID, "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorCompletedRequest) HTTPCode() int {
	return http.StatusBadRequest
}
