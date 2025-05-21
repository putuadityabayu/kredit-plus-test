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
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"xyz/internal/dto"
	"xyz/internal/model"
	"xyz/internal/service"
	"xyz/pkg/encrypt"
	"xyz/pkg/response"
	"xyz/pkg/validator"
)

func TestAuth_Login(t *testing.T) {
	mock := setupApp(t)
	svc := service.NewAuthService(mock.userRepo)
	defer mock.ctrl.Finish()

	validate := validator.New()
	password := "password"
	user := model.User{
		NIK:       "1234567890123456",
		FullName:  "User Test",
		LegalName: "User Legal Test",
		Salary:    6000000,
	}
	user.HashPassword(password)

	cases := []struct {
		name  string
		setup func() (req dto.LoginRequest, res *dto.LoginResponse, err error)
	}{
		{
			name: "Missing nik and password",
			setup: func() (req dto.LoginRequest, res *dto.LoginResponse, err error) {
				req = dto.LoginRequest{}

				err = validate.Struct(&req)
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
				return
			},
		},
		{
			name: "Missing nik",
			setup: func() (req dto.LoginRequest, res *dto.LoginResponse, err error) {
				req = dto.LoginRequest{
					Password: "password",
				}

				err = validate.Struct(&req)
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
				return
			},
		},
		{
			name: "Missing password",
			setup: func() (req dto.LoginRequest, res *dto.LoginResponse, err error) {
				req = dto.LoginRequest{
					NIK: "1234567890123456",
				}

				err = validate.Struct(&req)
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
				return
			},
		},
		{
			name: "Invalid NIK",
			setup: func() (req dto.LoginRequest, res *dto.LoginResponse, err error) {
				req = dto.LoginRequest{
					NIK:      "1234262362325266743734436346",
					Password: password,
				}

				err = validate.Struct(&req)
				err = response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
				return
			},
		},
		{
			name: "Get user error",
			setup: func() (req dto.LoginRequest, res *dto.LoginResponse, err error) {
				req = dto.LoginRequest{
					NIK:      user.NIK,
					Password: password,
				}

				err = errors.New("server error")
				mock.userRepo.EXPECT().GetByNIK(gomock.Any(), req.NIK).Return(nil, err).Times(1)
				err = response.ErrorServer("Internal server error", err)

				return
			},
		},
		{
			name: "User not found",
			setup: func() (req dto.LoginRequest, res *dto.LoginResponse, err error) {
				req = dto.LoginRequest{
					NIK:      user.NIK,
					Password: password,
				}

				mock.userRepo.EXPECT().GetByNIK(gomock.Any(), req.NIK).Return(nil, gorm.ErrRecordNotFound).Times(1)
				err = response.NotfoundHelper(gorm.ErrRecordNotFound, "Invalid nik or password")

				return
			},
		},
		{
			name: "Invalid password",
			setup: func() (req dto.LoginRequest, res *dto.LoginResponse, err error) {
				req = dto.LoginRequest{
					NIK:      user.NIK,
					Password: "wrong_password",
				}

				mock.userRepo.EXPECT().GetByNIK(gomock.Any(), req.NIK).Return(&user, nil).Times(1)

				err = response.ErrorParameter(response.ErrBadRequest, "Invalid nik or password", nil)
				return
			},
		},
		{
			name: "Generate JWT Error",
			setup: func() (req dto.LoginRequest, res *dto.LoginResponse, err error) {
				req = dto.LoginRequest{
					NIK:      user.NIK,
					Password: password,
				}

				mock.userRepo.EXPECT().GetByNIK(gomock.Any(), req.NIK).Return(&user, nil).Times(1)
				errJwt := errors.New("generate token error")

				monkey.Patch(encrypt.GenerateJWTToken, func(claims jwt.RegisteredClaims) (string, error) {
					return "", errJwt
				})

				err = response.ErrorServer("Failed to generate token", errJwt)
				return
			},
		},
		{
			name: "Success",
			setup: func() (req dto.LoginRequest, res *dto.LoginResponse, err error) {
				req = dto.LoginRequest{
					NIK:      user.NIK,
					Password: password,
				}

				mock.userRepo.EXPECT().GetByNIK(gomock.Any(), req.NIK).Return(&user, nil).Times(1)

				monkey.Patch(encrypt.GenerateJWTToken, func(claims jwt.RegisteredClaims) (string, error) {
					return "JWT Token", nil
				})

				res = &dto.LoginResponse{
					Token: "JWT Token",
					User:  user,
				}

				return
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, resExpected, errExpected := c.setup()
			defer monkey.UnpatchAll()

			res, err := svc.Login(context.Background(), req)
			if errExpected == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, err != nil)
				assert.Equal(t, errExpected, err)
				/*var e response.ErrorResponse
				if errors.As(err, &e) {
					assert.Equal(t, e.Code, errExpected.Code)

					assert.Equal(t, errExpected.Details, e.Details)
				} else {
					assert.ErrorContains(t, err, errExpected.Error())
				}*/
			}

			if resExpected != nil {
				assert.Equal(t, resExpected, res)
			}
		})
	}
}
