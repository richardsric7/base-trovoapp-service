package middleware

import "github.com/gin-gonic/gin"

func ExtractAddress(c *gin.Context) string {

	publicKey := c.GetHeader("X-TW-PUBLIC-KEY")

	return publicKey
}

func ExtractSigner(c *gin.Context) string {

	signer := c.GetHeader("X-TW-SIGNER")

	return signer
}

func ExtractSignature(c *gin.Context) string {

	sig := c.GetHeader("X-TW-SIGNATURE")

	return sig
}
func ExtractTimestamp(c *gin.Context) string {

	ts := c.GetHeader("X-TW-TIMESTAMP")

	return ts
}
func ExtractDeviceID(c *gin.Context) string {

	deviceID := c.GetHeader("X-TW-DEVICE-ID")

	return deviceID
}
