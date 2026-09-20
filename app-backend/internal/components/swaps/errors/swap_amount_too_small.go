package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorSwapAmountTooSmall = "error-swap-amount-too-small"

//ErrorSwapAmountTooSmall for an invalid payment amount
type ErrorSwapAmountTooSmall struct {
}

//Error returns the error string
func (e *ErrorSwapAmountTooSmall) Error() string {
	return errorSwapAmountTooSmall
}

//ErrorType returns error type as string
func (e *ErrorSwapAmountTooSmall) ErrorType() string {
	return errorSwapAmountTooSmall
}

//Data returns data of the error
func (e *ErrorSwapAmountTooSmall) Data() string {
	return "destinationAmount"
}

//Message returns tring message of error
func (e *ErrorSwapAmountTooSmall) Message() string {
	return "Swap amount too small"
}

//JSONError returns json of the error
func (e *ErrorSwapAmountTooSmall) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorSwapAmountTooSmall) HTTPCode() int {
	return http.StatusBadRequest
}
