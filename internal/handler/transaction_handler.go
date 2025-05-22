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

type TransactionHandler struct {
	transactionSvc service.TransactionService
}

func NewTransactionHandler(repo repository.RepoRegistry) TransactionHandler {
	userSvc := service.NewTransactionService(repo.UserRepository, repo.TransactionRepository)
	return TransactionHandler{
		transactionSvc: userSvc,
	}
}

func (h TransactionHandler) Create(c *fiber.Ctx) error {
	ctx, span := otel.StartSpan(c.UserContext(), "TransactionHandler.Create")
	defer span.End()
	c.SetUserContext(ctx)

	var req dto.TransactionRequest

	if err := c.BodyParser(&req); err != nil {
		span.RecordErrorHelper(response.ErrorServer("", err), "body parser")
		return response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	user, _, err := h.transactionSvc.Create(ctx, req)
	if err != nil {
		return err
	}

	return response.Success(c, user, fiber.StatusCreated, "Transaction created successfully")
}
