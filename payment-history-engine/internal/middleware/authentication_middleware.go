package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"trovo-wallet-payment-history-engine/internal/errors"
	"trovo-wallet-payment-history-engine/internal/validators"

	"github.com/gin-gonic/gin"
)

func authenticationChecks(keyParam string, c *gin.Context) error {
	keyParam = strings.TrimSpace(keyParam)
	fullUri := c.Request.URL.RequestURI()

	signerAddress := ExtractSigner(c)
	signature := ExtractSignature(c)
	log.Printf("Full Path With Query:[%s] KeyParam:[%s] Signature: [%s]\n", fullUri, keyParam, signature)

	publicKeyFormatError := validators.ValidateAddressFormat(signerAddress)

	if publicKeyFormatError != nil {
		return &errors.ErrorInvalidAuthenticationAddress{Address: signerAddress}
	}

	err := VerifyHttpSignature(fullUri, keyParam, signature, signerAddress)

	if err != nil {
		return err
	}

	return nil

}

func WebSocketAuthenticationChecks(body, signature, signerAddress string) error {

	publicKeyFormatError := validators.ValidateAddressFormat(signerAddress)

	if publicKeyFormatError != nil {
		return &errors.ErrorInvalidAuthenticationAddress{Address: signerAddress}
	}

	err := VerifyHttpSignature("", body, signature, signerAddress)

	if err != nil {
		return err
	}

	return nil

}

func AuthenticationMiddlewareUsingBody() gin.HandlerFunc {

	return func(c *gin.Context) {
		if os.Getenv("ENABLE_AUTH_MIDDLEWARE") == "0" {
			c.Next()
			return
		}
		h := c.Request.Header.Get("User-Agent")
		publicKey := ExtractAddress(c)

		log.Printf("[%s] is using [%s]\n", publicKey, h)

		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		// log.Println("body is:", string(body))
		authenticationError := authenticationChecks(string(body), c)

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
func AuthenticationMiddlewareUsingTimestamp() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("ENABLE_AUTH_MIDDLEWARE") == "0" {
			c.Next()
			return
		}
		h := c.Request.Header.Get("User-Agent")
		timestamp := ExtractTimestamp(c)
		signerAddress := ExtractSigner(c)

		log.Printf("[%s] is using [%s]\n", signerAddress, h)
		log.Printf("Timestamp:[%s] signerAddress:[%s]\n", timestamp, signerAddress)

		authenticationError := authenticationChecks(signerAddress+timestamp, c)

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
