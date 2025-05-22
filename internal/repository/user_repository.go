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

	Create(ctx context.Context, user *model.User, opts ...Option) error
	GetByID(ctx context.Context, id string, opts ...Option) (*model.User, error)
	GetByNIK(ctx context.Context, nik string, opts ...Option) (*model.User, error)
	Save(ctx context.Context, user *model.User, opts ...Option) error
	ListTenorLimits(ctx context.Context, userid string, opts ...Option) ([]*model.TenorLimits, error)
	ListTransactions(ctx context.Context, userid string, opts ...Option) (total int64, transactions []*model.Transaction, err error)
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

func (r userRepositoryImpl) Create(ctx context.Context, user *model.User, opts ...Option) error {
	return r.getDatabase(ctx, opts...).Create(user).Error
}

func (r userRepositoryImpl) GetByID(ctx context.Context, id string, opts ...Option) (*model.User, error) {
	var user model.User
	if err := r.getDatabase(ctx, opts...).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepositoryImpl) GetByNIK(ctx context.Context, nik string, opts ...Option) (*model.User, error) {
	var user model.User
	if err := r.getDatabase(ctx, opts...).Where("nik = ?", nik).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepositoryImpl) Save(ctx context.Context, user *model.User, opts ...Option) error {
	user.UpdatedAt = time.Now()
	return r.getDatabase(ctx, opts...).Save(user).Error
}

func (r userRepositoryImpl) ListTenorLimits(ctx context.Context, userid string, opts ...Option) ([]*model.TenorLimits, error) {
	var tenorLimits []*model.TenorLimits
	if err := r.getDatabase(ctx, opts...).Where("user_id = ?", userid).Order("tenor_in_months asc").Find(&tenorLimits).Error; err != nil {
		return nil, err
	}

	return tenorLimits, nil
}

func (r userRepositoryImpl) ListTransactions(ctx context.Context, userid string, opts ...Option) (total int64, transactions []*model.Transaction, err error) {
	db := r.getDatabase(ctx, opts...).Model(&model.Transaction{}).Where("user_id = ?", userid)

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Find(&transactions).Error
	if err != nil {
		return
	}

	return
}
