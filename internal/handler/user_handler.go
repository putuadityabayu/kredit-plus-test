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

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(repo repository.RepoRegistry) UserHandler {
	userSvc := service.NewUserService(repo.UserRepository)
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
	ctx, span := otel.StartSpan(c.UserContext(), "UserHandler.Update")
	defer span.End()
	c.SetUserContext(ctx)

	var req dto.UserRequest

	if err := c.BodyParser(&req); err != nil {
		span.RecordErrorHelper(response.ErrorServer("", err), "body parser")
		return response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	user, err := h.userSvc.Update(ctx, req)
	if err != nil {
		return err
	}

	return response.Success(c, user, fiber.StatusOK, "User updated successfully")
}

func (h UserHandler) ListNIK(c *fiber.Ctx) error {
	ctx, span := otel.StartSpan(c.UserContext(), "UserHandler.ListNIK")
	defer span.End()
	c.SetUserContext(ctx)

	limits, err := h.userSvc.GetTenorLimits(ctx)
	if err != nil {
		return err
	}

	return response.Success(c, limits, fiber.StatusOK, "Tenor limits retrieved successfully")
}

func (h UserHandler) ListTransactions(c *fiber.Ctx) error {
	ctx, span := otel.StartSpan(c.UserContext(), "UserHandler.ListTransactions")
	defer span.End()
	c.SetUserContext(ctx)

	var req dto.Pagination
	if err := c.QueryParser(&req); err != nil {
		span.RecordErrorHelper(response.ErrorServer("", err), "query parser")
		return response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	limits, meta, err := h.userSvc.GetTransactions(ctx, req)
	if err != nil {
		return err
	}

	return response.Success(c, limits, meta, fiber.StatusOK, "Tenor limits retrieved successfully")
}
