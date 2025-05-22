/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package dto

type TransactionRequest struct {
	OTR       float64 `json:"otr" validate:"required"`
	AssetName string  `json:"asset_name" validate:"required"`
	Tenor     int     `json:"tenor" validate:"required,min=1,max=6"` // Tenor yang dipilih (1, 2, 3, atau 6 bulan)
}
