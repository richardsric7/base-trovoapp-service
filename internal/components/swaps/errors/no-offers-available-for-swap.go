package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorSwapOfferNotAvailable = "error-swap-offer-not-available"

//ErrorSwapOfferNotAvailable for an invalid payment amount
type ErrorSwapOfferNotAvailable struct {
}

//Error returns the error string
func (e *ErrorSwapOfferNotAvailable) Error() string {
	return errorSwapOfferNotAvailable
}

//ErrorType returns error type as string
func (e *ErrorSwapOfferNotAvailable) ErrorType() string {
	return errorSwapOfferNotAvailable
}

//Data returns data of the error
func (e *ErrorSwapOfferNotAvailable) Data() string {
	return "destinationAsset"
}

//Message returns tring message of error
func (e *ErrorSwapOfferNotAvailable) Message() string {
	return "cannot swap at this time. try later"
}

//JSONError returns json of the error
func (e *ErrorSwapOfferNotAvailable) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

//HTTPCode returns http status code
func (e *ErrorSwapOfferNotAvailable) HTTPCode() int {
	return http.StatusBadRequest
}
