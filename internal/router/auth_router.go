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

func AuthRouterV1(app *fiber.App, repo repository.RepoRegistry) {
	routerV1 := app.Group("/v1")
	h := handler.NewAuthHandler(repo)

	routerV1.Post("/auth/login", h.Login)
}
