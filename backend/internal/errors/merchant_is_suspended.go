package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorServiceIsSuspended = "error-service-is-suspended"

//ErrorMerchantIsSuspended is returned during registration when a username already exists
type ErrorServiceIsSuspended struct {
	Username string
}

//Error returns the error string
func (e *ErrorServiceIsSuspended) Error() string {
	return errorServiceIsSuspended
}

//ErrorType returns error type as string
func (e *ErrorServiceIsSuspended) ErrorType() string {
	return errorServiceIsSuspended
}

//Data returns data of the error
func (e *ErrorServiceIsSuspended) Data() string {
	return "serviceLink"
}

//Message returns tring message of error
func (e *ErrorServiceIsSuspended) Message() string {
	return "Service Account is suspended"
}

//JSONError returns json of the error
func (e *ErrorServiceIsSuspended) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorServiceIsSuspended) HTTPCode() int {
	return http.StatusForbidden
}
