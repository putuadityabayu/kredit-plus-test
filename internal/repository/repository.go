/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package repository

import (
	"context"
	"xyz/internal/model"
)

type BaseRepository interface {
	StartTransaction(ctx context.Context, fc func(ctx context.Context) error) error
}

type UserRepository interface {
	BaseRepository

	Create(ctx context.Context, user *model.User) error
	GetByNIK(ctx context.Context, nik string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Save(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}
