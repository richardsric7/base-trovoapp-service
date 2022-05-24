package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorMerchantIsSuspended = "error-merchant-is-suspended"

//ErrorMerchantIsSuspended is returned during registration when a username already exists
type ErrorMerchantIsSuspended struct {
	Username string
}

//Error returns the error string
func (e *ErrorMerchantIsSuspended) Error() string {
	return errorMerchantIsSuspended
}

//ErrorType returns error type as string
func (e *ErrorMerchantIsSuspended) ErrorType() string {
	return errorMerchantIsSuspended
}

//Data returns data of the error
func (e *ErrorMerchantIsSuspended) Data() string {
	return "merchant"
}

//Message returns tring message of error
func (e *ErrorMerchantIsSuspended) Message() string {
	return "Merchant Account is suspended"
}

//JSONError returns json of the error
func (e *ErrorMerchantIsSuspended) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorMerchantIsSuspended) HTTPCode() int {
	return http.StatusForbidden
}
