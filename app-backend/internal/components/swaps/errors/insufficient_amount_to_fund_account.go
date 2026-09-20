package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInsufficientAmountToFundAccount = "error-insufficient-amount-to-fund-account"

//ErrorInsufficientAmountToFundAccount for an invalid payment amount
type ErrorInsufficientAmountToFundAccount struct {
}

//Error returns the error string
func (e *ErrorInsufficientAmountToFundAccount) Error() string {
	return errorInsufficientAmountToFundAccount
}

//ErrorType returns error type as string
func (e *ErrorInsufficientAmountToFundAccount) ErrorType() string {
	return errorInsufficientAmountToFundAccount
}

//Data returns data of the error
func (e *ErrorInsufficientAmountToFundAccount) Data() string {
	return "amount"
}

//Message returns tring message of error
func (e *ErrorInsufficientAmountToFundAccount) Message() string {
	return "Insufficient amount to fund new account. You need to send at least 3 XBN"
}

//JSONError returns json of the error
func (e *ErrorInsufficientAmountToFundAccount) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorInsufficientAmountToFundAccount) HTTPCode() int {
	return http.StatusBadRequest
}
