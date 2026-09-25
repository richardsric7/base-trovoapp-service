package middleware

import (
	"errors"
	"net/http"
)

// ExtractTokenMetadata extracts metaData from the request
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims := token.Claims
	if token.Valid {
		accessUUID := claims.AccessUUID

		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     claims.UserID,
		}, nil
	}
	return nil, err
}

// ExtractTokenMetadata extracts metaData from the request
func WebSocetExtractTokenMetadata(r string) (*AccessDetails, error) {
	token, err := WebsocketVerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims := token.Claims
	if token.Valid {
		accessUUID := claims.AccessUUID

		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     claims.UserID,
		}, nil
	}
	return nil, errors.New("invalid token")
}
