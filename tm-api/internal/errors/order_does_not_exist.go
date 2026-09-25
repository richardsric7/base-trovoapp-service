package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorOrderDoesNotExist = "error-order-does-not-exist"

// ErrorOrderDoesNotExist is when Order does not exist
type ErrorOrderDoesNotExist struct {
	OrderID string
}

// Error returns the error string
func (e *ErrorOrderDoesNotExist) Error() string {
	return errorOrderDoesNotExist
}

// ErrorType returns error type as string
func (e *ErrorOrderDoesNotExist) ErrorType() string {
	return errorOrderDoesNotExist
}

// Data returns data of the error
func (e *ErrorOrderDoesNotExist) Data() string {
	return fmt.Sprintf("%v", e.OrderID)
}

// Message returns tring message of error
func (e *ErrorOrderDoesNotExist) Message() string {
	return fmt.Sprintf("orderID: %v does not exist", e.OrderID)
}

// JSONError returns json of the error
func (e *ErrorOrderDoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorOrderDoesNotExist) HTTPCode() int {
	return http.StatusNotFound
}
