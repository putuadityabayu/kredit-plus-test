/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package encrypt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

const issuer = "xyz.com"

// GenerateJWTToken for generating token JWT for login
func GenerateJWTToken(claims jwt.RegisteredClaims) (string, error) {
	now := jwt.NewNumericDate(time.Now())
	claims.Issuer = issuer
	claims.IssuedAt = now
	claims.NotBefore = now

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(viper.GetString("secret.jwt")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWTToken(tokenString string) (*jwt.RegisteredClaims, error) {
	// parse and validating token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(viper.GetString("secret.jwt")), nil
	}, jwt.WithIssuer(issuer), jwt.WithExpirationRequired())

	if err != nil {
		return nil, err
	}

	// check token if it is valid or not
	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	var claims jwt.RegisteredClaims
	byt, _ := json.Marshal(claimsMap)
	_ = json.Unmarshal(byt, &claims)

	return &claims, nil
}
