/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

import "github.com/gofiber/fiber/v2"

func ErrorServer(msg string, err error) ErrorResponse {
	return NewError(fiber.StatusInternalServerError, "SERVER_ERROR", msg, nil, err)
}
