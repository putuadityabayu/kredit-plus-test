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

type User struct {
	ID             string    `gorm:";column:id;primaryKey;type:uuid" json:"id"`
	NIK            string    `json:"nik" gorm:";column:nik;unique;type:varchar(16)"`
	FullName       string    `json:"full_name" gorm:"type:varchar(255)"`
	LegalName      string    `json:"legal_name" gorm:"column:legal_name;type:varchar(255)"`
	BirthPlace     string    `json:"birth_place"  gorm:"column:birth_place;type:varchar(255)"`
	BirthDateDb    string    `json:"-" gorm:"column:birth_date;type:date"`
	BirthDate      time.Time `json:"birth_date" gorm:"-"`
	Salary         float64   `json:"salary" gorm:"column:salary;type:decimal"`
	KTPPhotoURL    string    `json:"ktp_photo_url" gorm:"column:ktp_photo_url;type:varchar(255)"`
	SelfiePhotoURL string    `json:"selfie_photo_url" gorm:"column:selfie_photo_url;type:varchar(255)"`
	Date

	Password string `json:"-" gorm:"column:password;type:varchar(255)"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeSave(_ *gorm.DB) error {
	u.BirthDateDb = u.BirthDate.Format("2006-01-02")
	return nil
}

func (u *User) AfterFind(_ *gorm.DB) error {
	u.BirthDate, _ = time.Parse("2006-01-02", u.BirthDateDb)
	return nil
}
