/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package service

import (
	"context"
	"xyz/internal/model"
	"xyz/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return UserServiceImpl{
		userRepository: userRepository,
	}
}

func (u UserServiceImpl) GetByID(ctx context.Context, id string) (*model.User, error) {
	return u.userRepository.GetByID(ctx, id)
}
