/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
	"xyz/internal/model"
)

type UserRepository interface {
	BaseRepository

	Create(ctx context.Context, user *model.User) error
	GetByNIK(ctx context.Context, nik string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Save(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}
type userRepositoryImpl struct {
	base
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		base: base{
			db: db,
		},
	}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	return r.getDatabase(ctx).Create(user).Error
}

func (r *userRepositoryImpl) GetByNIK(ctx context.Context, nik string) (*model.User, error) {
	var user model.User
	if err := r.getDatabase(ctx).Where("nik = ?", nik).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.getDatabase(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	return r.getDatabase(ctx).Save(user).Error
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.getDatabase(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}
