package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorTransactionMismatch = "error-transaction-mismatch"

// ErrorTransactionMismatch for when generated transaction does not match supplied transaction
type ErrorTransactionMismatch struct {
}

// Error returns the error string
func (e *ErrorTransactionMismatch) Error() string {
	return errorTransactionMismatch
}

// ErrorType returns error type as string
func (e *ErrorTransactionMismatch) ErrorType() string {
	return errorTransactionMismatch
}

// Data returns data of the error
func (e *ErrorTransactionMismatch) Data() string {
	return "transaction"
}

// Message returns tring message of error
func (e *ErrorTransactionMismatch) Message() string {
	return "Expected transaction does not match supplied transaction"
}

// JSONError returns json of the error
func (e *ErrorTransactionMismatch) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorTransactionMismatch) HTTPCode() int {
	return http.StatusBadRequest
}
