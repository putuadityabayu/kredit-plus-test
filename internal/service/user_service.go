/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package service

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.portalnesia.com/utils"
	"time"
	"xyz/internal/dto"
	"xyz/internal/model"
	"xyz/internal/repository"
	"xyz/pkg/helper"
	"xyz/pkg/otel"
	"xyz/pkg/response"
	"xyz/pkg/validator"
)

type UserService interface {
	Create(ctx context.Context, user dto.UserRequest) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user dto.UserRequest) (*model.User, error)
	GetTenorLimits(ctx context.Context) ([]*model.TenorLimits, error)
	GetTransactions(ctx context.Context, req dto.Pagination) ([]*model.Transaction, *response.Meta, error)
}

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return userServiceImpl{
		userRepository: userRepository,
	}
}

func (u userServiceImpl) Create(ctx context.Context, req dto.UserRequest) (*model.User, error) {
	var span *otel.Span
	ctx, span = otel.StartSpan(ctx, "UserService.Create")
	defer span.End()

	validate := validator.New()

	// validate request with validator
	if err := validate.Struct(req); err != nil {
		span.RecordErrorHelper(err, "validator")
		return nil, response.ErrorParameter(response.ErrBadRequest, response.MsgInvalidRequest, err)
	}

	// next validation
	errs := response.NewErrorFields()

	// validate birthday
	_, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		errs.Add("birth_date", "Invalid date format")
	}

	// password required
	if req.Password == "" {
		errs.Add("password", "Password is required")
	}

	// confirm password required
	if req.ConfirmPassword == "" {
		errs.Add("confirm_password", "Confirm password is required")
	}

	// password and confirm password must be the same
	if req.Password != req.ConfirmPassword {
		errs.Add("password", "Password and confirm password must be the same")
	}

	// return error if there is any
	if errs.Exist() {
		span.RecordErrorHelper(errs, "validation")
		return nil, response.ErrorParameter(response.ErrBadRequest, response.MsgInvalidRequest, errs)
	}

	user := &model.User{
		ID:         utils.UUID(),
		NIK:        req.NIK,
		FullName:   req.FullName,
		LegalName:  req.LegalName,
		BirthPlace: req.BirthPlace,
		BirthDate:  req.BirthDate,
		Salary:     req.Salary,
		Date:       model.NewDate(),
	}
	// hash password
	user.HashPassword(req.Password)

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		span.RecordErrorHelper(err, "Create data error")
		return nil, response.DatabaseHelper(err, map[string]string{"idx_users_nik": "NIK"}, span)
	}

	return user, nil
}

func (u userServiceImpl) GetByID(ctx context.Context, id string) (*model.User, error) {
	var span *otel.Span
	ctx, span = otel.StartSpan(ctx, "UserService.GetByID")
	defer span.End()

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		span.RecordErrorHelper(err, "repository.GetByID")
		return nil, response.NotfoundHelper(err, "user not found", span)
	}

	return user, nil
}

func (u userServiceImpl) Update(ctx context.Context, req dto.UserRequest) (*model.User, error) {
	var span *otel.Span
	ctx, span = otel.StartSpan(ctx, "UserService.Update")
	defer span.End()

	userid := helper.GetValueContext(ctx, "userid", "")
	if userid == "" {
		return nil, response.Authorization(fiber.StatusUnauthorized, response.ErrUnauthorized, response.MsgLoginRequired)
	}

	validate := validator.New()

	// validate request with validator
	if err := validate.Struct(req); err != nil {
		span.RecordErrorHelper(err, "validator")
		return nil, response.ErrorParameter(response.ErrBadRequest, response.MsgInvalidRequest, err)
	}

	// next validation
	errs := response.NewErrorFields()

	// validate birthday
	_, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		errs.Add("birth_date", "Invalid date format")
	}

	// return error if there is any
	if errs.Exist() {
		span.RecordErrorHelper(err, "validation")
		return nil, response.ErrorParameter(response.ErrBadRequest, response.MsgInvalidRequest, errs)
	}

	var user *model.User
	err = u.userRepository.StartTransaction(ctx, func(ctx context.Context) error {
		var errTx error

		// get user by id
		user, errTx = u.userRepository.GetByID(ctx, userid)
		if errTx != nil {
			return response.NotfoundHelper(errTx, "user not found", span)
		}

		user.FullName = req.FullName
		user.LegalName = req.LegalName
		user.BirthPlace = req.BirthPlace
		user.BirthDate = req.BirthDate
		user.Salary = req.Salary

		// update user
		errTx = u.userRepository.Save(ctx, user)
		if errTx != nil {
			return response.DatabaseHelper(errTx, map[string]string{"idx_user_nik": "NIK"}, span)
		}

		return nil
	})
	if err != nil {
		span.RecordErrorHelper(err, "db.transaction")
		return nil, err
	}

	return user, nil
}

func (u userServiceImpl) GetTenorLimits(ctx context.Context) ([]*model.TenorLimits, error) {
	var span *otel.Span
	ctx, span = otel.StartSpan(ctx, "UserService.GetTenorLimits")
	defer span.End()

	userid := helper.GetValueContext(ctx, "userid", "")
	if userid == "" {
		return nil, response.Authorization(fiber.StatusUnauthorized, response.ErrUnauthorized, response.MsgLoginRequired)
	}

	tenorLimits, err := u.userRepository.ListTenorLimits(ctx, userid)
	if err != nil {
		span.RecordErrorHelper(err, "repository.ListTenorLimits")
		return nil, response.ErrorServer(response.MsgInternalServer, err)
	}

	return tenorLimits, nil
}

func (u userServiceImpl) GetTransactions(ctx context.Context, req dto.Pagination) ([]*model.Transaction, *response.Meta, error) {
	var span *otel.Span
	ctx, span = otel.StartSpan(ctx, "UserService.GetTransactions")
	defer span.End()

	userid := helper.GetValueContext(ctx, "userid", "")
	if userid == "" {
		return nil, nil, response.Authorization(fiber.StatusUnauthorized, response.ErrUnauthorized, response.MsgLoginRequired)
	}

	total, transactions, err := u.userRepository.ListTransactions(ctx, userid, repository.WithPagination(req))
	if err != nil {
		span.RecordErrorHelper(err, "repository.ListTransactions")
		return nil, nil, response.ErrorServer(response.MsgInternalServer, err)
	}

	meta := req.Meta(total)

	return transactions, &meta, nil
}
