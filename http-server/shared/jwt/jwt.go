package jwtHelper

import (
	"crypto/rsa"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type JwtClaim struct {
	Account string
	jwt.StandardClaims
}

func ValidateToken(signedToken string, publicKey *rsa.PublicKey) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return
}
