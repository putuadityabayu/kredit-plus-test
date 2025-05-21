/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"xyz/internal/dto"
	"xyz/internal/repository"
	"xyz/pkg/encrypt"
	"xyz/pkg/response"
	"xyz/pkg/validator"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type AuthServiceImpl struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return AuthServiceImpl{
		userRepository: userRepository,
	}
}

func (s AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	validate := validator.New()

	// validate request with validator
	if err := validate.Struct(req); err != nil {
		return nil, response.ErrorParameter(response.ErrBadRequest, "Invalid request parameter", err)
	}

	// get user by nik
	user, err := s.userRepository.GetByNIK(ctx, req.NIK)
	if err != nil {
		return nil, response.NotfoundHelper(err, "Invalid nik or password")
	}

	// check user password
	validPassword := user.CheckPassword(req.Password)
	if !validPassword {
		return nil, response.ErrorParameter(response.ErrBadRequest, "Invalid nik or password", nil)
	}

	// generate jwt
	token, err := encrypt.GenerateJWTToken(jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	})
	if err != nil {
		return nil, response.ErrorServer("Failed to generate token", err)
	}

	resp := &dto.LoginResponse{
		Token: token,
		User:  *user,
	}

	return resp, nil
}
