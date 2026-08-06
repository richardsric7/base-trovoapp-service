package users

import (
	"fmt"
	"log"
	"net/http"
	"os"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/golang-jwt/jwt"
)

func LogUserIn(user userModels.ServiceLinksUser, gc *sharedconfig.GlobalConfig) (interface{}, error) {
	//create tokens
	//Create new pairs of refresh and access tokens
	// cacheKey := "[GET] /v1/users"
	ts, createErr := middleware.CreateToken(user.ID)
	if createErr != nil {
		log.Printf("[LogUserIn] error creating new token: %v\n", createErr)

		return nil, createErr
	}
	//save the tokens metadata to redis
	saveErr := middleware.CreateAuth(user.ID, ts, gc.RedisCache)
	if saveErr != nil {
		log.Printf("[LogUserIn] error saving new token: %v\n", saveErr)

		// c.JSON(http.StatusForbidden, gin.H{"error": saveErr.Error(), "message": saveErr.Error()})
		return nil, saveErr
	}

	tokenResponse := struct {
		AccessToken  string                         `json:"accessToken"`
		RefreshToken string                         `json:"refreshToken"`
		UserInfo     userModels.ServiceLinkUserInfo `json:"userInfo"`
	}{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
		UserInfo:     user.ToServiceLinkUserInfo(gc),
	}

	return tokenResponse, nil
}

func RefreshToken(refreshToken string, gc *sharedconfig.GlobalConfig) (interface{}, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		log.Printf("[RefreshToken] err token parse: %v\n", err)

		return nil, &tErrors.CustomError{Param: "token", Err: "error-refresh-token-expired",
			ErrMessage: "refresh token expired",
			Code:       http.StatusUnauthorized}
	}

	if !token.Valid {
		return nil, &tErrors.CustomError{Param: "token", Err: "error-refresh-token-expired",
			ErrMessage: "refresh token expired",
			Code:       http.StatusUnauthorized}
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return nil, &tErrors.CustomError{Param: "token", Err: "error-refresh-token-expired",
				ErrMessage: "refresh token expired",
				Code:       http.StatusUnauthorized}
		}
		userID := claims["user_id"].(string)

		//Delete the previous Refresh Token
		deleted, delErr := middleware.DeleteAuth(refreshUUID, gc.RedisCache)
		if delErr != nil || deleted == 0 { //if any goes wrong
			log.Printf("[RefreshToken] error deleting old refresh token: %v, %v\n", delErr, deleted)
			return nil, delErr
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := middleware.CreateToken(userID)
		if createErr != nil {
			log.Printf("[RefreshToken] error creating new token: %v\n", createErr)

			return nil, createErr
		}
		//save the tokens metadata to redis
		saveErr := middleware.CreateAuth(userID, ts, gc.RedisCache)
		if saveErr != nil {
			log.Printf("[RefreshToken] error saving new token: %v\n", saveErr)
			return nil, saveErr
		}

		u, _ := usersDB.GetUser(userID, gc.DB, gc)

		ui := u.ToServiceLinkUser(gc)

		tokenResponse := struct {
			AccessToken  string                         `json:"accessToken"`
			RefreshToken string                         `json:"refreshToken"`
			UserInfo     userModels.ServiceLinkUserInfo `json:"userInfo"`
		}{
			AccessToken:  ts.AccessToken,
			RefreshToken: ts.RefreshToken,
			UserInfo:     ui.ToServiceLinkUserInfo(gc),
		}
		return tokenResponse, nil
	} else {
		return nil, &tErrors.CustomError{Param: "token", Err: "error-refresh-token-expired",
			ErrMessage: "refresh token expired",
			Code:       http.StatusUnauthorized}
	}

}
