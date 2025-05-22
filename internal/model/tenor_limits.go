/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package model

import "time"

type TenorLimits struct {
	ID            string  `gorm:";column:id;primaryKey;type:uuid" json:"id"`
	UserID        string  `gorm:";column:user_id;type:uuid" json:"-"`
	TenorInMonths int     `gorm:";column:tenor_in_months;type:int" json:"tenor_in_months"`
	LimitAmount   float64 `gorm:";column:limit_amount;type:int" json:"limit_amount"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;autoUpdateTime"`

	User *User `json:"user,omitempty" gorm:"<-:false;foreignKey:user_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (t *TenorLimits) TableName() string {
	return "user_tenor_limits"
}
