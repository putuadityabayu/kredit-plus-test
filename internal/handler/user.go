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
	"xyz/pkg/helper"
	"xyz/pkg/response"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userRepo repository.UserRepository) UserHandler {
	userSvc := service.NewUserService(userRepo)
	return UserHandler{
		userSvc: userSvc,
	}
}

func (h UserHandler) Create(c *fiber.Ctx) error {
	var req dto.UserRequest

	if err := c.BodyParser(&req); err != nil {
		return response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	user, err := h.userSvc.Create(c.UserContext(), req)
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
	userid := helper.GetValueContext(c.UserContext(), "userid", "")

	if userid == "" {
		return response.Authorization(fiber.StatusForbidden, "FORBIDDEN", "You don't have permission to access this resource")
	}

	var req dto.UserRequest

	if err := c.BodyParser(&req); err != nil {
		return response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	user, err := h.userSvc.Update(c.UserContext(), userid, req)
	if err != nil {
		return err
	}

	return response.Success(c, user, fiber.StatusOK, "User updated successfully")
}
