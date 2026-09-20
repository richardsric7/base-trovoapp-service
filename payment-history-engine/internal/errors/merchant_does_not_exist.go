package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorMerchantDoesNotExist = "error-merchant-does-not-exist"

// ErrorMerchantDoesNotExist is when user does not exist
type ErrorMerchantDoesNotExist struct {
	Username string
}

// Error returns the error string
func (e *ErrorMerchantDoesNotExist) Error() string {
	return errorMerchantDoesNotExist
}

// ErrorType returns error type as string
func (e *ErrorMerchantDoesNotExist) ErrorType() string {
	return errorMerchantDoesNotExist
}

// Data returns data of the error
func (e *ErrorMerchantDoesNotExist) Data() string {
	return fmt.Sprintf("%v", e.Username)
}

// Message returns tring message of error
func (e *ErrorMerchantDoesNotExist) Message() string {
	return fmt.Sprintf("merchant: %v does not exist", e.Username)
}

// JSONError returns json of the error
func (e *ErrorMerchantDoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorMerchantDoesNotExist) HTTPCode() int {
	return http.StatusNotFound
}
