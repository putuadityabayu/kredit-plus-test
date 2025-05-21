/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package model

import (
	"github.com/spf13/viper"
	pncrypto "go.portalnesia.com/crypto"
	"go.portalnesia.com/nullable"
)

type User struct {
	ID             string          `gorm:";column:id;primaryKey;type:uuid" json:"id"`
	NIK            nullable.String `json:"nik" gorm:";column:nik;unique;type:varchar(16)"`
	FullName       string          `json:"full_name" gorm:"type:varchar(255)"`
	LegalName      nullable.String `json:"legal_name" gorm:"column:legal_name;type:varchar(255)"`
	BirthPlace     nullable.String `json:"birth_place"  gorm:"column:birth_place;type:varchar(255)"`
	BirthDate      nullable.String `json:"birth_date" gorm:"column:birth_date;type:date"`
	Salary         nullable.Float  `json:"salary" gorm:"column:salary;type:decimal"`
	KTPPhotoURL    nullable.String `json:"ktp_photo_url" gorm:"column:ktp_photo_url;type:varchar(255)"`
	SelfiePhotoURL nullable.String `json:"selfie_photo_url" gorm:"column:selfie_photo_url;type:varchar(255)"`
	Date

	Password string `json:"-" gorm:"column:password;type:varchar(255)"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) HashPassword(passwordString string) {
	saltPassword := passwordString + viper.GetString("secret.password_salt")
	hashPassword := pncrypto.HashPassword(saltPassword)
	u.Password = hashPassword
}

func (u *User) CheckPassword(passwordString string) bool {
	saltPassword := passwordString + viper.GetString("secret.password_salt")
	return pncrypto.ComparePassword(saltPassword, u.Password)
}
