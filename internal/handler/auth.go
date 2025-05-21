/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package handler

import (
	"github.com/gofiber/fiber/v2"
	"xyz/internal/dto"
	"xyz/internal/repository"
	"xyz/internal/service"
	"xyz/pkg/otel"
	"xyz/pkg/response"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(userRepo repository.UserRepository) AuthHandler {
	authSvc := service.NewAuthService(userRepo)
	return AuthHandler{
		authSvc: authSvc,
	}
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	ctx, span := otel.StartSpan(c.UserContext(), "AuthHandler.Login")
	defer span.End()
	c.SetUserContext(ctx)

	if err := c.BodyParser(&req); err != nil {
		span.RecordErrorHelper(response.ErrorServer("", err), "body parser")
		return response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	user, err := h.authSvc.Login(ctx, req)
	if err != nil {
		return err
	}

	return response.Success(c, user, fiber.StatusOK, "Login success")
}
