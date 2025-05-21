/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package service

import (
	"context"
	"go.portalnesia.com/nullable"
	"go.portalnesia.com/utils"
	"time"
	"xyz/internal/dto"
	"xyz/internal/model"
	"xyz/internal/repository"
	"xyz/pkg/response"
	"xyz/pkg/validator"
)

type UserService interface {
	Create(ctx context.Context, user dto.UserRequest) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, id string, user dto.UserRequest) (*model.User, error)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return UserServiceImpl{
		userRepository: userRepository,
	}
}

func (u UserServiceImpl) Create(ctx context.Context, req dto.UserRequest) (*model.User, error) {
	validate := validator.New()

	// validate request with validator
	if err := validate.Struct(req); err != nil {
		return nil, response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	// next validation
	errs := make([]any, 0)

	// validate birthday
	_, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		errs = append(errs, response.NewErrorFields(
			[2]string{"birth_date", "Invalid date format"},
		))
	}

	// password required
	if req.Password == "" {
		errs = append(errs, response.NewErrorFields(
			[2]string{"password", "Password is required"},
		))
	}

	// confirm password required
	if req.ConfirmPassword == "" {
		errs = append(errs, response.NewErrorFields(
			[2]string{"confirm_password", "Confirm password is required"},
		))
	}

	// password and confirm password must be the same
	if req.Password != req.ConfirmPassword {
		errs = append(errs, response.NewErrorFields(
			[2]string{"password", "Password and confirm password must be the same"},
		))
	}

	// return error if there is any
	if len(errs) > 0 {
		return nil, response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", errs...)
	}

	user := &model.User{
		ID:         utils.UUID(),
		NIK:        nullable.NewString(req.NIK),
		FullName:   req.FullName,
		LegalName:  nullable.NewString(req.LegalName),
		BirthPlace: nullable.NewString(req.BirthPlace),
		BirthDate:  nullable.NewString(req.BirthDate),
		Salary:     nullable.NewFloat(req.Salary),
		Date:       model.NewDate(),
	}
	// hash password
	user.HashPassword(req.Password)

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u UserServiceImpl) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, response.NotfoundHelper(err, "user not found")
	}

	return user, nil
}

func (u UserServiceImpl) Update(ctx context.Context, id string, req dto.UserRequest) (*model.User, error) {
	validate := validator.New()

	// validate request with validator
	if err := validate.Struct(req); err != nil {
		return nil, response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	// next validation
	errs := make([]any, 0)

	// validate birthday
	_, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		errs = append(errs, response.NewErrorFields(
			[2]string{"birth_date", "Invalid date format"},
		))
	}

	// return error if there is any
	if len(errs) > 0 {
		return nil, response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", errs...)
	}

	var user *model.User
	err = u.userRepository.StartTransaction(ctx, func(ctx context.Context) error {
		var errTx error

		// get user by id
		user, errTx = u.userRepository.GetByID(ctx, id)
		if errTx != nil {
			return errTx
		}

		user.FullName = req.FullName
		user.LegalName = nullable.NewString(req.LegalName)
		user.BirthPlace = nullable.NewString(req.BirthPlace)
		user.BirthDate = nullable.NewString(req.BirthDate)
		user.Salary = nullable.NewFloat(req.Salary)

		// update user
		errTx = u.userRepository.Save(ctx, user)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
