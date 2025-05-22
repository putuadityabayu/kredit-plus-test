/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package test

import (
	"bou.ke/monkey"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.portalnesia.com/utils"
	"gorm.io/gorm"
	"testing"
	"xyz/internal/dto"
	"xyz/internal/model"
	"xyz/internal/service"
	"xyz/pkg/response"
	"xyz/pkg/validator"
)

func TestUser_Create(t *testing.T) {
	mock := setupApp(t)
	svc := service.NewUserService(mock.userRepo, mock.limitRepo)
	defer mock.ctrl.Finish()

	validate := validator.New()
	password := "password123"
	tmpReq := dto.UserRequest{
		NIK:             "1234567890",
		FullName:        "John Doe",
		LegalName:       "John Doe",
		BirthPlace:      "New York",
		BirthDate:       "1990-01-01",
		Salary:          5000000,
		Password:        password,
		ConfirmPassword: password,
	}

	cases := []struct {
		name  string
		setup func() (req dto.UserRequest, res *model.User, err error)
	}{
		{
			name: "Missing required fields",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = dto.UserRequest{}
				err = validate.Struct(&req)
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
				return
			},
		},
		{
			name: "Invalid date format",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq
				req.BirthDate = "invalid-date"

				errs := response.NewErrorFields()
				errs.Add("birth_date", "Invalid date format")
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", errs)
				return
			},
		},
		{
			name: "Password missing",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq
				req.Password = ""
				errs := response.NewErrorFields()
				errs.Add("password", "Password is required")
				errs.Add("password", "Password and confirm password must be the same")
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", errs)
				return
			},
		},
		{
			name: "Confirm password missing",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq
				req.ConfirmPassword = ""
				errs := response.NewErrorFields()
				errs.Add("confirm_password", "Confirm password is required")
				errs.Add("password", "Password and confirm password must be the same")
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", errs)
				return
			},
		},
		{
			name: "Password mismatch",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq
				req.Password = "another-password"
				errs := response.NewErrorFields()
				errs.Add("password", "Password and confirm password must be the same")
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", errs)
				return
			},
		},
		{
			name: "Successful creation",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq

				monkey.Patch(utils.UUID, func() string {
					return "test-id"
				})

				res = &model.User{
					ID:         "test-id",
					NIK:        req.NIK,
					FullName:   req.FullName,
					LegalName:  req.LegalName,
					BirthPlace: req.BirthPlace,
					BirthDate:  req.BirthDate,
					Salary:     req.Salary,
					Password:   req.Password,
				}
				res.HashPassword(password)

				mock.userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				return
			},
		},
		{
			name: "Repository error",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq
				errs := errors.New("repository error")
				err = response.ErrorServer("Internal server error", errs)
				mock.userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errs).Times(1)
				return
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, resExpected, expectedErr := tc.setup()
			res, err := svc.Create(context.Background(), req)
			defer monkey.UnpatchAll()

			if expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, err != nil)
				assert.Equal(t, expectedErr, err)
			}

			if resExpected != nil {
				assert.NotNil(t, res)
				// bypass password
				res.Password = ""
				resExpected.Password = ""
				//bypass date
				date := model.NewDate()
				res.Date = date
				resExpected.Date = date

				assert.Equal(t, resExpected, res)
			}
		})
	}
}

func TestUser_GetByID(t *testing.T) {
	mock := setupApp(t)
	svc := service.NewUserService(mock.userRepo, mock.limitRepo)
	defer mock.ctrl.Finish()

	cases := []struct {
		name  string
		setup func() (id string, res *model.User, err error)
	}{
		{
			name: "User found",
			setup: func() (id string, res *model.User, err error) {
				id = "test-id"
				res = &model.User{
					ID:       "test-id",
					NIK:      "1234567890",
					FullName: "John Doe",
				}
				mock.userRepo.EXPECT().GetByID(gomock.Any(), id).Return(res, nil)
				return
			},
		},
		{
			name: "User not found",
			setup: func() (id string, res *model.User, err error) {
				id = "non-existent-id"

				errs := gorm.ErrRecordNotFound
				err = response.NotfoundHelper(errs, "user not found")
				mock.userRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errs)
				return
			},
		},
		{
			name: "Repository error",
			setup: func() (id string, res *model.User, err error) {
				id = "test-id"

				errs := errors.New("database error")
				err = response.NotfoundHelper(errs, "user not found")
				mock.userRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errs)
				return
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, resExpected, expectedErr := tc.setup()
			res, err := svc.GetByID(context.Background(), id)

			if expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, err != nil)
				assert.Equal(t, expectedErr, err)
			}

			if resExpected != nil {
				assert.NotNil(t, res)
				// bypass password
				res.Password = ""
				resExpected.Password = ""
				//bypass date
				date := model.NewDate()
				res.Date = date
				resExpected.Date = date

				assert.Equal(t, resExpected, res)
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	mock := setupApp(t)
	svc := service.NewUserService(mock.userRepo, mock.limitRepo)
	defer mock.ctrl.Finish()

	validate := validator.New()
	id := "test-id"
	tmpReq := dto.UserRequest{
		NIK:        "1234567890",
		FullName:   "John Doe",
		LegalName:  "John Doe",
		BirthPlace: "New York",
		BirthDate:  "1990-01-01",
		Salary:     5000000,
	}

	cases := []struct {
		name  string
		setup func() (req dto.UserRequest, res *model.User, err error)
	}{
		{
			name: "Missing required fields",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = dto.UserRequest{}
				err = validate.Struct(&req)
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
				return
			},
		},
		{
			name: "Invalid date format",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq
				req.BirthDate = "invalid-date"

				errs := response.NewErrorFields()
				errs.Add("birth_date", "Invalid date format")
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", errs)
				return
			},
		},
		{
			name: "User not found",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq

				errs := gorm.ErrRecordNotFound
				err = response.NotfoundHelper(errs, "user not found")
				mock.userRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errs).Times(1)
				return
			},
		},
		{
			name: "Repository error",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq
				resp := &model.User{
					ID:  id,
					NIK: req.NIK,
				}
				mock.userRepo.EXPECT().GetByID(gomock.Any(), id).Return(resp, nil)
				errs := errors.New("repository error")
				mock.userRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errs).Times(1)
				err = response.ErrorServer("Internal server error", errs)
				return
			},
		},
		{
			name: "Successful update",
			setup: func() (req dto.UserRequest, res *model.User, err error) {
				req = tmpReq
				res = &model.User{
					ID:         id,
					NIK:        req.NIK,
					FullName:   req.FullName,
					LegalName:  req.LegalName,
					BirthPlace: req.BirthPlace,
					BirthDate:  req.BirthDate,
					Salary:     req.Salary,
				}
				mock.userRepo.EXPECT().GetByID(gomock.Any(), id).Return(res, nil)
				mock.userRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				return
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, resExpected, expectedErr := tc.setup()
			res, err := svc.Update(context.Background(), id, req)

			if expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, err != nil)
				assert.Equal(t, expectedErr, err)
			}

			if resExpected != nil {
				assert.NotNil(t, res)
				// bypass password
				res.Password = ""
				resExpected.Password = ""
				//bypass date
				date := model.NewDate()
				res.Date = date
				resExpected.Date = date

				assert.Equal(t, resExpected, res)
			}
		})
	}
}
