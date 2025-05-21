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

type UserRepositoryImpl struct {
	base
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		base: base{
			db: db,
		},
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	return r.getDatabase(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) GetByNIK(ctx context.Context, nik string) (*model.User, error) {
	var user model.User
	if err := r.getDatabase(ctx).Where("nik = ?", nik).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.getDatabase(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	return r.getDatabase(ctx).Save(user).Error
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.getDatabase(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}
