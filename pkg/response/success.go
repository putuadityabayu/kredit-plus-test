/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

import "github.com/gofiber/fiber/v2"

type Meta struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
}

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Success response
//
// opts:
// - meta: Meta
// - http_status: int
// - message: string
func Success(c *fiber.Ctx, data interface{}, opts ...any) error {
	var (
		status  = 200
		meta    *Meta
		message = ""
	)

	if len(opts) > 0 {
		for _, o := range opts {
			switch v := o.(type) {
			case int:
				status = v
			case *Meta:
				meta = v
			case string:
				message = v
			}
		}
	}

	if status == 204 {
		return c.SendStatus(status)
	}

	resp := SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	return c.Status(status).JSON(resp)
}
