/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package dto

import (
	"math"
	"xyz/pkg/response"
)

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"per_page"`
}

func (p Pagination) Meta(totalItems int64) response.Meta {
	return response.Meta{
		TotalItems:  totalItems,
		TotalPages:  int(math.Ceil(float64(totalItems) / float64(p.Limit))),
		CurrentPage: p.Page,
		PerPage:     p.Limit,
	}
}
