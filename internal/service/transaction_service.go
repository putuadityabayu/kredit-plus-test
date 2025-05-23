/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package service

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.portalnesia.com/utils"
	"gorm.io/gorm"
	"time"
	"xyz/internal/dto"
	"xyz/internal/model"
	"xyz/internal/repository"
	"xyz/pkg/helper"
	"xyz/pkg/otel"
	"xyz/pkg/response"
	"xyz/pkg/validator"
)

type TransactionService interface {
	Create(ctx context.Context, req dto.TransactionRequest) (*model.Transaction, *model.TenorLimits, error)
}

type transactionServiceImpl struct {
	userRepository        repository.UserRepository
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(userRepository repository.UserRepository, transactionRepository repository.TransactionRepository) TransactionService {
	return transactionServiceImpl{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
	}
}

func (t transactionServiceImpl) Create(ctx context.Context, req dto.TransactionRequest) (*model.Transaction, *model.TenorLimits, error) {
	var span *otel.Span
	ctx, span = otel.StartSpan(ctx, "TransactionService.Create")
	defer span.End()

	userid := helper.GetValueContext(ctx, "userid", "")
	if userid == "" {
		return nil, nil, response.Authorization(fiber.StatusUnauthorized, response.ErrUnauthorized, response.MsgLoginRequired)
	}

	validate := validator.New()

	// validate request with validator
	if err := validate.Struct(req); err != nil {
		span.RecordErrorHelper(err, "validator")
		return nil, nil, response.ErrorParameter(response.ErrBadRequest, response.MsgInvalidRequest, err)
	}

	var (
		trx   *model.Transaction
		limit *model.TenorLimits
	)
	err := t.transactionRepository.StartTransaction(ctx, func(ctx context.Context) error {
		// get user
		user, err := t.userRepository.GetByID(ctx, userid)
		if err != nil {
			return response.NotfoundHelper(err, "User not found", span)
		}

		date := time.Now()
		// create transaction
		trx = &model.Transaction{
			ID:              utils.UUID(),
			ContractNumber:  helper.GenerateContractNumber(),
			UserID:          user.ID,
			OTR:             req.OTR,
			AssetName:       req.AssetName,
			Tenor:           req.Tenor,
			TransactionDate: date,
			Status:          model.TrxPENDING,
			CreatedAt:       date,
			UpdatedAt:       date,
		}

		var totalAmount float64
		trx.InterestAmount, trx.AdminFee, trx.InstallmentAmount, totalAmount = helper.GetTransactionAmount(req.OTR, req.Tenor)

		// get limit
		limit, err = t.transactionRepository.GetLimit(ctx, user.ID, req.Tenor, repository.WithLockTable())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ErrorParameter(response.ErrInsufficientLimit, response.MsgInsufficientLimit, fiber.StatusUnprocessableEntity)
			}
			return response.ErrorServer(response.MsgInternalServer, err)
		}

		// check limit
		if limit.LimitAmount < totalAmount {
			return response.ErrorParameter(response.ErrInsufficientLimit, response.MsgInsufficientLimit, fiber.StatusUnprocessableEntity)
		}
		limit.LimitAmount -= totalAmount

		// save transaction
		if err = t.transactionRepository.Create(ctx, trx); err != nil {
			return response.ErrorServer(response.MsgInternalServer, err)
		}

		// update limit
		if err = t.transactionRepository.UpdateTenorLimit(ctx, limit); err != nil {
			return response.ErrorServer(response.MsgInternalServer, err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return trx, limit, nil
}
