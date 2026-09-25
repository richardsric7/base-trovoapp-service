package middleware

import (
	"admin-panel-dashboard/internal/trovosdk"
	"net/http"
)

// DeleteAuth delete the UUID from redis server
func DeleteAuth(r *http.Request) error {
	tokenString := ExtractToken(r)
	sl, err := trovosdk.NewServiceLink()
	if err != nil {

		return err
	}
	_, err = sl.JwtTokenDelete(tokenString)
	if err != nil {

		return err
	}

	return nil
}
