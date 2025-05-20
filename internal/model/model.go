/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package model

import (
	"gorm.io/gorm"
	"time"
)

type Date struct {
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp;index"`
}

func NewDate() Date {
	now := time.Now()
	return Date{
		CreatedAt: now,
		UpdatedAt: &now,
	}
}
