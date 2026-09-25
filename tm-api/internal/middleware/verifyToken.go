package middleware

import (
	"admin-panel-dashboard/internal/trovosdk"
	"errors"
	"net/http"
)

// // VerifyToken verifies the JWT token from the header
func VerifyToken(r *http.Request) (*trovosdk.Token, error) {
	tokenString := ExtractToken(r)
	sl, err := trovosdk.NewServiceLink()
	if err != nil {
		// log.Println("[JwtTokenAuthMiddleware] error:", err)

		return &trovosdk.Token{}, err
	}
	tokenResponse, err := sl.JwtTokenVerify(tokenString)
	if err != nil {

		return &trovosdk.Token{}, err
	}
	token := tokenResponse.Token
	if !token.Valid {
		return token, errors.New("error invalid token")
	}
	return token, nil
}

// VerifyToken verifies the JWT token from the header
func WebsocketVerifyToken(r string) (*trovosdk.Token, error) {
	sl, err := trovosdk.NewServiceLink()
	if err != nil {
		// log.Println("[JwtTokenAuthMiddleware] error:", err)

		return &trovosdk.Token{}, err
	}
	tokenResponse, err := sl.JwtTokenVerify(r)
	if err != nil {

		return &trovosdk.Token{}, err
	}
	if !tokenResponse.Token.Valid {
		return &trovosdk.Token{}, errors.New("invalid token")
	}
	return tokenResponse.Token, nil
}
