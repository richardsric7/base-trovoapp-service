package middleware

import (
	"log"
	"net/http"
	"os"
	"trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Merchant holds Merchant data model
type Merchant struct {
	ID        string `json:"-" gorm:"size:100"`
	ApiKey    string `json:"apiKey" gorm:"size:50;index:idx_merchant_api_key,unique;not null;"`
	Suspended int    `json:"-" gorm:"type:integer;not null;default:0"`
}

func apiKeyChecks(c *gin.Context, gc *sharedconfig.GlobalConfig) error {
	// apiKey = strings.TrimSpace(apiKey)
	fullUri := c.Request.URL.RequestURI()
	merchantKey := ExtractMerchantApiKey(c)
	log.Printf("Full Path With Query:[%s] APIKEY:[%s]\n", fullUri, merchantKey)

	if len(merchantKey) == 0 {
		return &errors.CustomError{Param: "apiKey", Err: "Error Missing APIKEY parameter", ErrMessage: "Missing APIKEY parameter"}
	}

	err := VerifyMerchantAPIKey(merchantKey, gc.DB)

	if err != nil {
		return err
	}

	return nil

}

// VerifyMerchantAPIKey checks merchant by API key
func VerifyMerchantAPIKey(apiKey string, db *gorm.DB) (err error) {

	e := db.Where(Merchant{ApiKey: apiKey}).First(&Merchant{}).Error

	if e != nil {

		return &errors.CustomError{Param: "apiKey", Err: "Error Invalid API Key", ErrMessage: "Invalid API Key", Code: http.StatusForbidden}

	}

	return nil

}
func AuthenticationMiddlewareUsingAPIKey(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("ENABLE_AUTH_MIDDLEWARE") == "0" {
			c.Next()
			return
		}
		h := c.Request.Header.Get("User-Agent")
		merchantKey := ExtractMerchantApiKey(c)

		log.Printf("[%s] is using [%s]\n", merchantKey, h)

		authenticationError := apiKeyChecks(c, gc)

		if authenticationError != nil {
			var ex errors.GenericError
			var ok bool

			ex, ok = authenticationError.(errors.GenericError)
			if ok {
				c.JSON(http.StatusUnauthorized, ex.JSONError())
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": authenticationError})
			}
			c.Abort()
			return
		}

		c.Next()
	}
}
