/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package test

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
	"xyz/internal/dto"
	"xyz/internal/model"
	"xyz/internal/service"
	"xyz/pkg/response"
	"xyz/pkg/validator"
)

func TestTransactionService_Create(t *testing.T) {
	mock := setupApp(t)
	svc := service.NewTransactionService(mock.userRepo, mock.transactionRepo)
	defer mock.ctrl.Finish()

	date := time.Now().Add(-24 * time.Hour)
	validate := validator.New()
	userId := "user-id"
	tmpReq := dto.TransactionRequest{
		OTR:       800000,
		AssetName: "Test asset",
		Tenor:     3,
	}

	tmpLimit := model.TenorLimits{
		ID:            "tenor-limit-id",
		UserID:        "user-id",
		TenorInMonths: 3,
		LimitAmount:   1000000,
		CreatedAt:     date,
		UpdatedAt:     date,
	}
	user := &model.User{
		ID:       "user-id",
		NIK:      "1234567890",
		FullName: "John Doe",
	}

	cases := []struct {
		name          string
		setup         func() (req dto.TransactionRequest, res *model.Transaction, err error)
		notLogin      bool
		expectedLimit float64
	}{
		{
			name: "User not logged in",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				err = response.Authorization(fiber.StatusUnauthorized, response.ErrUnauthorized, response.MsgLoginRequired)
				return
			},
			notLogin: true,
		},
		{
			name: "Invalid request",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				req = dto.TransactionRequest{}
				err = validate.Struct(&req)
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
				return
			},
		},
		{
			name: "User not found",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				req = tmpReq
				errs := gorm.ErrRecordNotFound
				err = response.NotfoundHelper(errs, "User not found")
				mock.userRepo.EXPECT().GetByID(gomock.Any(), userId).Return(nil, errs)
				return
			},
		},
		{
			name: "Get user error",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				req = tmpReq
				errs := errors.New("database error get user")
				err = response.NotfoundHelper(errs, "user not found")
				mock.userRepo.EXPECT().GetByID(gomock.Any(), userId).Return(nil, errs)
				return
			},
		},
		{
			name: "Get limit error",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				req = tmpReq

				mock.userRepo.EXPECT().GetByID(gomock.Any(), userId).Return(user, nil)

				errs := errors.New("database error get limit")
				err = response.ErrorServer(response.MsgInternalServer, errs)
				mock.transactionRepo.EXPECT().GetLimit(gomock.Any(), userId, req.Tenor, gomock.Any()).Return(nil, errs)
				return
			},
		},
		{
			name: "Get limit not found",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				req = tmpReq

				mock.userRepo.EXPECT().GetByID(gomock.Any(), userId).Return(user, nil)

				err = response.ErrorParameter(response.ErrInsufficientLimit, response.MsgInsufficientLimit, fiber.StatusUnprocessableEntity)
				mock.transactionRepo.EXPECT().GetLimit(gomock.Any(), userId, req.Tenor, gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
				return
			},
		},
		{
			name: "Credit limit exceeded",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				limit := tmpLimit
				limit.TenorInMonths = 1
				limit.LimitAmount = 100000

				req = tmpReq
				req.OTR = 100000
				req.Tenor = 1

				mock.userRepo.EXPECT().GetByID(gomock.Any(), userId).Return(user, nil)
				mock.transactionRepo.EXPECT().GetLimit(gomock.Any(), userId, req.Tenor, gomock.Any()).Return(&limit, nil)
				err = response.ErrorParameter(response.ErrInsufficientLimit, response.MsgInsufficientLimit, fiber.StatusUnprocessableEntity)
				return
			},
		},
		{
			name: "Save transaction error",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				limit := tmpLimit
				req = tmpReq

				mock.userRepo.EXPECT().GetByID(gomock.Any(), userId).Return(user, nil)
				mock.transactionRepo.EXPECT().GetLimit(gomock.Any(), userId, req.Tenor, gomock.Any()).Return(&limit, nil)

				errs := errors.New("database error get limit")
				err = response.ErrorServer(response.MsgInternalServer, errs)
				mock.transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errs)

				return
			},
		},
		{
			name: "Save limit error",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				limit := tmpLimit
				req = tmpReq

				mock.userRepo.EXPECT().GetByID(gomock.Any(), userId).Return(user, nil)
				mock.transactionRepo.EXPECT().GetLimit(gomock.Any(), userId, req.Tenor, gomock.Any()).Return(&limit, nil)
				mock.transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

				errs := errors.New("database error save limit")
				err = response.ErrorServer(response.MsgInternalServer, errs)
				mock.transactionRepo.EXPECT().UpdateTenorLimit(gomock.Any(), gomock.Any()).Return(errs)

				return
			},
		},
		{
			name: "Create transaction success",
			setup: func() (req dto.TransactionRequest, res *model.Transaction, err error) {
				limit := tmpLimit
				req = tmpReq

				mock.userRepo.EXPECT().GetByID(gomock.Any(), userId).Return(user, nil)
				mock.transactionRepo.EXPECT().GetLimit(gomock.Any(), userId, req.Tenor, gomock.Any()).Return(&limit, nil)
				mock.transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				mock.transactionRepo.EXPECT().UpdateTenorLimit(gomock.Any(), gomock.Any()).Return(nil)

				return
			},
			expectedLimit: 200000,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			if !c.notLogin {
				ctx = context.WithValue(ctx, "userid", userId)
			}

			req, expectedRes, expectedErr := c.setup()

			res, limit, err := svc.Create(ctx, req)

			if expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, err != nil)
				assert.Equal(t, expectedErr, err)
			}

			if expectedRes != nil {
				assert.NotNil(t, res)

				assert.Equal(t, expectedRes, res)

				if c.expectedLimit != 0 {
					assert.Equal(t, c.expectedLimit, limit.LimitAmount)
				}
			}
		})
	}
}
