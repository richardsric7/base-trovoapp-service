package middleware

import "github.com/gin-gonic/gin"

func ExtractPublicKey(c *gin.Context) string {

	publicKey := c.GetHeader("X-TW-PUBLIC-KEY")

	return publicKey
}

func ExtractSigner(c *gin.Context) string {

	sig := c.GetHeader("X-TW-SIGNER")

	return sig
}

func ExtractSignature(c *gin.Context) string {

	sig := c.GetHeader("X-TW-SIGNATURE")

	return sig
}
func ExtractTimestamp(c *gin.Context) string {

	sig := c.GetHeader("X-TW-TIMESTAMP")

	return sig
}
func ExtractDeviceID(c *gin.Context) string {

	sig := c.GetHeader("X-TW-DEVICE-ID")

	return sig
}



