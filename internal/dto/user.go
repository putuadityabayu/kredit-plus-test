/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package dto

type UserRequest struct {
	NIK             string  `json:"nik" validate:"required,max=16"`
	FullName        string  `json:"full_name" validate:"required"`
	LegalName       string  `json:"legal_name" validate:"required"`
	BirthPlace      string  `json:"birth_place" validate:"required"`
	BirthDate       string  `json:"birth_date" validate:"required"`
	Salary          float64 `json:"salary" validate:"required"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirm_password"`
}
