package jwtHelper

import (
	"crypto/rsa"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type JwtClaim struct {
	Account string
	jwt.StandardClaims
}

func GenerateToken(account string, expireTime int64, secretKey *rsa.PrivateKey) (signedToken string, err error) {
	claims := &JwtClaim{
		Account: account,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(expireTime)).Unix(),
			Issuer:    "TCP",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err = token.SignedString(secretKey)

	if err != nil {
		return
	}

	return
}
