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
	"xyz/internal/middleware"
	"xyz/internal/repository"
)

func UserRouterV1(app *fiber.App, repo repository.RepoRegistry) {
	routerV1 := app.Group("/v1")
	h := handler.NewUserHandler(repo)

	routerV1.Get("/user/tenor", middleware.AuthorizationCheck, h.ListNIK)
	routerV1.Get("/user/:id", middleware.AuthorizationCheck, h.GetByID)
	routerV1.Post("/user", h.Create)
	routerV1.Put("/user", middleware.Authorization, h.Update)
}
