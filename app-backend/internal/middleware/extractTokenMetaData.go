package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// ExtractTokenMetadata extracts metaData from the request
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     claims["user_id"].(string),
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
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     claims["user_id"].(string),
		}, nil
	}
	return nil, err
}
