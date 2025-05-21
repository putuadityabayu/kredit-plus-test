/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package dto

import "xyz/internal/model"

type LoginRequest struct {
	NIK      string `json:"nik" validate:"required,max=16"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}
