/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package handler

import (
	"github.com/gofiber/fiber/v2"
	"xyz/internal/service"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) UserHandler {
	return UserHandler{
		userSvc: userSvc,
	}
}

func (h UserHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := h.userSvc.GetByID(c.UserContext(), id)
	if err != nil {
		return err
	}

	return c.Next()
}
