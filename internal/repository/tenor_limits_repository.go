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

type TenorLimitsRepository interface {
	BaseRepository

	ListByUserID(ctx context.Context, userid string) ([]model.TenorLimits, error)
}

type tenorLimitsRepository struct {
	base
}

func NewTenorLimitsRepository(db *gorm.DB) TenorLimitsRepository {
	return &tenorLimitsRepository{
		base: base{
			db: db,
		},
	}
}

func (r *tenorLimitsRepository) ListByUserID(ctx context.Context, userid string) ([]model.TenorLimits, error) {
	var tenorLimits []model.TenorLimits
	if err := r.getDatabase(ctx).Where("user_id = ?", userid).Order("tenor_in_months asc").Find(&tenorLimits).Error; err != nil {
		return nil, err
	}

	return tenorLimits, nil
}
