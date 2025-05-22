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

func TransactionRouterV1(app *fiber.App, repo repository.RepoRegistry) {
	routerV1 := app.Group("/v1")
	h := handler.NewTransactionHandler(repo)

	routerV1.Post("/transaction", middleware.Authorization, h.Create)
}
