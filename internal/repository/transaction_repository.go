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
	"xyz/internal/model"
)

type TransactionRepository interface {
	BaseRepository

	Create(ctx context.Context, user *model.Transaction, opts ...Option) error
	GetLimit(ctx context.Context, userId string, tenor int, opts ...Option) (*model.TenorLimits, error)
	UpdateTenorLimit(ctx context.Context, tenorLimit *model.TenorLimits, opts ...Option) error
}
type transactionRepositoryImpl struct {
	base
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{
		base: base{
			db: db,
		},
	}
}

func (r transactionRepositoryImpl) Create(ctx context.Context, transaction *model.Transaction, opts ...Option) error {
	return r.getDatabase(ctx, opts...).Create(transaction).Error
}

func (r transactionRepositoryImpl) GetLimit(ctx context.Context, userId string, tenor int, opts ...Option) (*model.TenorLimits, error) {
	var tenorLimit model.TenorLimits
	if err := r.getDatabase(ctx, opts...).Where("user_id = ? AND tenor_in_months = ?", userId, tenor).First(&tenorLimit).Error; err != nil {
		return nil, err
	}
	return &tenorLimit, nil
}

func (r transactionRepositoryImpl) UpdateTenorLimit(ctx context.Context, tenorLimit *model.TenorLimits, opts ...Option) error {
	return r.getDatabase(ctx, opts...).Save(tenorLimit).Error
}
