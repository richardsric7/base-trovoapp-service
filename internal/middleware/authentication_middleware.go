package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/validators"

	"github.com/gin-gonic/gin"
)

func authenticationChecks(keyParam string, c *gin.Context) error {
	keyParam = strings.TrimSpace(keyParam)
	fullUri := c.Request.URL.RequestURI()
	// log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^Full URI:", fullUri)
	log.Printf("Full Path and query is [%s]\n", fullUri)

	publicKey := ExtractPublicKey(c)
	signature := ExtractSignature(c)

	publicKeyFormatError := validators.ValidatePublicKeyFormat(publicKey)

	if publicKeyFormatError != nil {
		return &errors.ErrorInvalidAuthenticationPublicKey{PublicKey: publicKey}
	}

	err := VerifyHttpSignature(fullUri, keyParam, signature, publicKey)

	if err != nil {
		return err
	}

	return nil

}

func WebSocketAuthenticationChecks(body, signature, signerPublicKey string) error {

	publicKeyFormatError := validators.ValidatePublicKeyFormat(signerPublicKey)

	if publicKeyFormatError != nil {
		return &errors.ErrorInvalidAuthenticationPublicKey{PublicKey: signerPublicKey}
	}

	err := VerifyHttpSignature("", body, signature, signerPublicKey)

	if err != nil {
		return err
	}

	return nil

}

func AuthenticationMiddlewareUsingBody() gin.HandlerFunc {

	return func(c *gin.Context) {
		if os.Getenv("ENABLE_AUTH_MIDDLEWARE") == "0" {
			c.Next()
		}
		h := c.Request.Header.Get("User-Agent")
		publicKey := ExtractPublicKey(c)

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
		}
		h := c.Request.Header.Get("User-Agent")
		timestamp := c.Request.Header.Get("X-TW-TIMESTAMP")
		publicKey := ExtractPublicKey(c)

		log.Printf("[%s] is using [%s]\n", publicKey, h)

		authenticationError := authenticationChecks(publicKey+timestamp, c)

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
