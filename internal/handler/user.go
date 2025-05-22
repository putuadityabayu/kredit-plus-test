/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"xyz/internal/dto"
	"xyz/internal/repository"
	"xyz/internal/service"
	"xyz/pkg/helper"
	"xyz/pkg/otel"
	"xyz/pkg/response"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(repo repository.RepoRegistry) UserHandler {
	userSvc := service.NewUserService(repo.UserRepository, repo.LimitRepository)
	return UserHandler{
		userSvc: userSvc,
	}
}

func (h UserHandler) Create(c *fiber.Ctx) error {
	ctx, span := otel.StartSpan(c.UserContext(), "UserHandler.Create")
	defer span.End()
	c.SetUserContext(ctx)

	var req dto.UserRequest

	if err := c.BodyParser(&req); err != nil {
		span.RecordErrorHelper(response.ErrorServer("", err), "body parser")
		return response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	user, err := h.userSvc.Create(ctx, req)
	if err != nil {
		return err
	}

	return response.Success(c, user, fiber.StatusCreated, "User created successfully")
}

func (h UserHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userSvc.GetByID(c.UserContext(), id)
	if err != nil {
		return err
	}

	return response.Success(c, user, fiber.StatusOK, "User retrieved successfully")
}

func (h UserHandler) Update(c *fiber.Ctx) error {
	ctx, span := otel.StartSpan(c.UserContext(), "UserHandler.Create")
	defer span.End()
	c.SetUserContext(ctx)

	userid := helper.GetValueContext(c.UserContext(), "userid", "")

	if userid == "" {
		span.RecordErrorHelper(response.ErrorServer("", errors.New("missing userid")), "userid == \"\"")
		return response.Authorization(fiber.StatusForbidden, "FORBIDDEN", "You don't have permission to access this resource")
	}

	var req dto.UserRequest

	if err := c.BodyParser(&req); err != nil {
		span.RecordErrorHelper(response.ErrorServer("", err), "body parser")
		return response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	user, err := h.userSvc.Update(ctx, userid, req)
	if err != nil {
		return err
	}

	return response.Success(c, user, fiber.StatusOK, "User updated successfully")
}

func (h UserHandler) ListNIK(c *fiber.Ctx) error {
	ctx, span := otel.StartSpan(c.UserContext(), "UserHandler.Create")
	defer span.End()
	c.SetUserContext(ctx)
	userid := helper.GetValueContext(c.UserContext(), "userid", "")

	if userid == "" {
		span.RecordErrorHelper(response.ErrorServer("", errors.New("missing userid")), "userid == \"\"")
		return response.Authorization(fiber.StatusForbidden, "FORBIDDEN", "You don't have permission to access this resource")
	}
	limits, err := h.userSvc.GetTenorLimits(ctx, userid)
	if err != nil {
		return err
	}

	return response.Success(c, limits, fiber.StatusOK, "Tenor limits retrieved successfully")
}
