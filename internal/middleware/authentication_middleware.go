package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/validators"

	bantupaysdk "github.com/bantublockchain/bantupaysdk-go/security"
	"github.com/gin-gonic/gin"
)

func authenticationChecks(body string, c *gin.Context) error {
	body = strings.TrimSpace(body)
	fullUri := c.Request.URL.RequestURI()
	// log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^Full URI:", fullUri)
	log.Printf("Full Path and query is [%s]\n", fullUri)

	publicKey := ExtractPublicKey(c)
	signature := ExtractSignature(c)

	publicKeyFormatError := validators.ValidatePublicKeyFormat(publicKey)

	if publicKeyFormatError != nil {
		return &errors.ErrorInvalidAuthenticationPublicKey{PublicKey: publicKey}
	}

	err := bantupaysdk.VerifyHttpSignature(fullUri, body, signature, publicKey)

	if err != nil {
		return err
	}

	return nil

}

func WebSocketAuthenticationChecks(body, signature, publicKey string) error {

	publicKeyFormatError := validators.ValidatePublicKeyFormat(publicKey)

	if publicKeyFormatError != nil {
		return &errors.ErrorInvalidAuthenticationPublicKey{PublicKey: publicKey}
	}

	err := bantupaysdk.VerifyHttpSignature("", body, signature, publicKey)

	if err != nil {
		return err
	}

	return nil

}

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

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
