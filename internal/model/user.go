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
	"gorm.io/gorm"
)

type User struct {
	ID             string          `gorm:";column:id;primaryKey;type:uuid" json:"id"`
	NIK            string          `json:"nik" gorm:";column:nik;unique;type:varchar(16)"`
	FullName       string          `json:"full_name" gorm:"type:varchar(255)"`
	LegalName      string          `json:"legal_name" gorm:"column:legal_name;type:varchar(255)"`
	BirthPlace     string          `json:"birth_place"  gorm:"column:birth_place;type:varchar(255)"`
	BirthDate      string          `json:"birth_date" gorm:"column:birth_date;type:date"`
	Salary         float64         `json:"salary" gorm:"column:salary;type:decimal"`
	KTPPhotoURL    nullable.String `json:"ktp_photo_url" gorm:"column:ktp_photo_url;type:varchar(255)"`
	SelfiePhotoURL nullable.String `json:"selfie_photo_url" gorm:"column:selfie_photo_url;type:varchar(255)"`
	Date

	Password string `json:"-" gorm:"column:password;type:varchar(255)"`

	TenorLimits []TenorLimits `json:"tenor_limits,omitempty" gorm:"<-:false;foreignKey:user_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) AfterDelete(tx *gorm.DB) error {
	return tx.Model(&User{}).Where("id=?", u.ID).Unscoped().Updates(map[string]string{
		"nik": u.NIK + "--deleted",
	}).Error
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
