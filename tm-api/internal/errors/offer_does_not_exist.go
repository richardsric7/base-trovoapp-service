package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const errorOfferDoesNotExist = "error-offer-does-not-exist"

// ErrorOfferDoesNotExist is when Offer does not exist
type ErrorOfferDoesNotExist struct {
	OfferID string
}

// Error returns the error string
func (e *ErrorOfferDoesNotExist) Error() string {
	return errorOfferDoesNotExist
}

// ErrorType returns error type as string
func (e *ErrorOfferDoesNotExist) ErrorType() string {
	return errorOfferDoesNotExist
}

// Data returns data of the error
func (e *ErrorOfferDoesNotExist) Data() string {
	return fmt.Sprintf("%v", e.OfferID)
}

// Message returns tring message of error
func (e *ErrorOfferDoesNotExist) Message() string {
	return fmt.Sprintf("offerID: %v does not exist", e.OfferID)
}

// JSONError returns json of the error
func (e *ErrorOfferDoesNotExist) JSONError() gin.H {
	return gin.H{"error": e.ErrorType(), "data": e.Data(), "message": e.Message()}
}

// HTTPCode returns http status code
func (e *ErrorOfferDoesNotExist) HTTPCode() int {
	return http.StatusNotFound
}
