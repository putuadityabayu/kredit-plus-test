/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"xyz/internal/dto"
)

type options struct {
	db *gorm.DB
}

type Option func(db *gorm.DB) *gorm.DB

// WithLockTable locks the table
func WithLockTable() Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.Locking{
			Strength: "UPDATE",
		})
	}
}

func WithPagination(req *dto.Pagination) Option {
	return func(db *gorm.DB) *gorm.DB {
		if req.Page == 0 {
			req.Page = 1
		}
		if req.Limit == 0 {
			req.Limit = 10
		}
		offset := (req.Page - 1) * req.Limit
		return db.Offset(offset).Limit(req.Limit)
	}
}
