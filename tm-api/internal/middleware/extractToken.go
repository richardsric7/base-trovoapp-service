package middleware

import (
	"net/http"
	"strings"
)

// ExtractToken extracts the JWT token from the header
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	// normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	} else if len(strArr) == 1 {
		return strArr[0]
	}
	return ""
}
