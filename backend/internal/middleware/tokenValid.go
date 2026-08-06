package middleware

import (
	"net/http"
)

// TokenValid checks validity of the token. returns nil for valid token
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if !token.Valid {
		return err
	}

	return nil
}

// TokenValid checks validity of the token. returns nil for valid token
func WebsocketTokenValid(r string) error {
	token, err := WebsocketVerifyToken(r)
	if err != nil {
		return err
	}
	if !token.Valid {
		return err
	}

	return nil
}
