/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package model

import (
	"go.portalnesia.com/nullable"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             string          `gorm:";column:id;primaryKey;type:uuid" json:"id"`
	NIK            nullable.String `json:"nik" gorm:";column:nik;unique;type:varchar(16)"`
	FullName       string          `json:"full_name" gorm:"type:varchar(255)"`
	LegalName      nullable.String `json:"legal_name" gorm:"column:legal_name;type:varchar(255)"`
	BirthPlace     nullable.String `json:"birth_place"  gorm:"column:birth_place;type:varchar(255)"`
	BirthDateDb    nullable.String `json:"-" gorm:"column:birth_date;type:date"`
	BirthDate      nullable.Time   `json:"birth_date" gorm:"-"`
	Salary         nullable.Float  `json:"salary" gorm:"column:salary;type:decimal"`
	KTPPhotoURL    nullable.String `json:"ktp_photo_url" gorm:"column:ktp_photo_url;type:varchar(255)"`
	SelfiePhotoURL nullable.String `json:"selfie_photo_url" gorm:"column:selfie_photo_url;type:varchar(255)"`
	Date

	Password string `json:"-" gorm:"column:password;type:varchar(255)"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeSave(_ *gorm.DB) error {
	if u.BirthDate.Valid {
		u.BirthDateDb = nullable.NewString(u.BirthDate.Data.Format("2006-01-02"))
	}
	return nil
}

func (u *User) AfterFind(_ *gorm.DB) error {
	if u.BirthDateDb.Valid == false {
		u.BirthDate = nullable.NewTime(time.Now())
		u.BirthDate.Data, _ = time.Parse("2006-01-02", u.BirthDateDb.Data)
	}
	return nil
}
