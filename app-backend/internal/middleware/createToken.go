package middleware

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TokenDetails Token details struct
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// CreateToken creates token for logged in users
func CreateToken(userid string) (*TokenDetails, error) {
	td := &TokenDetails{}
	var err error
	te := time.Duration(10)
	tr := time.Duration(15)
	rte, err := decimal.NewFromString(os.Getenv("JWT_TOKEN_EXPIRY"))
	if err == nil && rte.IsPositive() {
		te = time.Duration(rte.IntPart())
	}
	rtr, err := decimal.NewFromString(os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"))
	if err == nil && rtr.IsPositive() {
		tr = time.Duration(rtr.IntPart())
	}
	td.AtExpires = time.Now().Add(time.Minute * te).Unix()
	td.AccessUUID = uuid.NewString()

	td.RtExpires = time.Now().Add(time.Minute * tr).Unix()
	td.RefreshUUID = uuid.NewString()

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}
