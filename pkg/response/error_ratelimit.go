/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

import "github.com/gofiber/fiber/v2"

func ErrorRateLimit() ErrorResponse {
	return NewError(fiber.StatusTooManyRequests, "RATE_LIMIT", "Too many requests", nil)
}
