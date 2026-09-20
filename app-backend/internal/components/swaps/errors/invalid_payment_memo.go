package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorInvalidPaymentMemo = "error-invalid-payment-memo"

//ErrorInvalidPaymentMemo for an invalid payment memo
type ErrorInvalidPaymentMemo struct {
}

//Error returns the error string
func (e *ErrorInvalidPaymentMemo) Error() string {
	return errorInvalidPaymentMemo
}

//ErrorType returns error type as string
func (e *ErrorInvalidPaymentMemo) ErrorType() string {
	return errorInvalidPaymentMemo
}

//Data returns data of the error
func (e *ErrorInvalidPaymentMemo) Data() string {
	return "memo"
}

//Message returns tring message of error
func (e *ErrorInvalidPaymentMemo) Message() string {
	return "Invalid payment memo. Memos cannot be longer than 28 characters"
}

//JSONError returns json of the error
func (e *ErrorInvalidPaymentMemo) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorInvalidPaymentMemo) HTTPCode() int {
	return http.StatusBadRequest
}
