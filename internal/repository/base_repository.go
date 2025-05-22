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
)

type BaseRepository interface {
	StartTransaction(ctx context.Context, fc func(ctx context.Context) error) error
}

type base struct {
	db *gorm.DB
}

func (b *base) StartTransaction(ctx context.Context, fc func(ctx context.Context) error) error {
	return b.db.Transaction(func(tx *gorm.DB) error {
		// save transaction to context
		ctx = context.WithValue(ctx, "db", tx)
		errTx := fc(ctx)
		if errTx != nil {
			return errTx
		}
		return nil
	})
}

func (b *base) getDatabase(ctx context.Context) *gorm.DB {
	dbAny := ctx.Value("db")
	if dbAny != nil {
		dbGorm, ok := dbAny.(*gorm.DB)
		if ok {
			return dbGorm
		}
	}
	return b.db
}
