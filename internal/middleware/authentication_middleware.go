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

	log.Printf("Full Path With Query:[%s] KeyParam:[%s]\n", fullUri, keyParam)

	signerPublicKey := ExtractSigner(c)
	signature := ExtractSignature(c)

	publicKeyFormatError := validators.ValidatePublicKeyFormat(signerPublicKey)

	if publicKeyFormatError != nil {
		return &errors.ErrorInvalidAuthenticationPublicKey{PublicKey: signerPublicKey}
	}

	err := VerifyHttpSignature(fullUri, keyParam, signature, signerPublicKey)

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
			return
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
			return
		}
		h := c.Request.Header.Get("User-Agent")
		timestamp := ExtractTimestamp(c)
		signerPublicKey := ExtractSigner(c)

		log.Printf("[%s] is using [%s]\n", signerPublicKey, h)
		log.Printf("Timestamp:[%s] signerPublicKey:[%s]\n", timestamp, signerPublicKey)

		authenticationError := authenticationChecks(signerPublicKey+timestamp, c)

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
