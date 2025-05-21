/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"strings"
	"xyz/pkg/encrypt"
	"xyz/pkg/response"
)

// Authorization is middleware to check authorization
func Authorization(c *fiber.Ctx) error {
	if err := authorization(c); err != nil {
		return err
	}

	return c.Next()
}

// AuthorizationCheck is middleware to check authorization only and not return error
func AuthorizationCheck(c *fiber.Ctx) error {
	_ = authorization(c)
	return c.Next()
}

func authorization(c *fiber.Ctx) error {
	authHeader := c.Get("authorization", "")

	if authHeader == "" {
		return response.Authorization(fiber.StatusUnauthorized, response.ErrUnauthorized, response.MsgMissingAuthorization)
	}

	authSplit := strings.Split(authHeader, " ")
	authType := strings.ToLower(authSplit[0])
	authToken := authSplit[1]

	if authType != "bearer" {
		return response.Authorization(fiber.StatusUnauthorized, response.ErrUnauthorized, response.MsgInvalidToken)
	}

	// parse and validating token
	claims, err := encrypt.ValidateJWTToken(authToken)
	if err != nil {
		return response.Authorization(fiber.StatusUnauthorized, response.ErrUnauthorized, err.Error())
	}

	ctx := context.WithValue(c.UserContext(), "userid", claims.Subject)
	c.SetUserContext(ctx)

	return nil
}
