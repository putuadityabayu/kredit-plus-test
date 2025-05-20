/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package router

import (
	"github.com/gofiber/fiber/v2"
	"xyz/internal/handler"
	"xyz/internal/repository"
)

func UserRouterV1(app *fiber.App, userRepo repository.UserRepository) {
	router := app.Group("/v1")
	h := handler.NewUserHandler(userRepo)

	router.Get("/users")
	router.Get("/user/:id", h.GetByID)
	router.Post("/user")
	router.Put("/user/:id")
	router.Delete("/user/:id")
}
