/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

import "github.com/gofiber/fiber/v2"

const (
	ErrNotFound = "NOT_FOUND"
	MsgNotFound = "%s with %s `%s` not found"
)

func NotFound(message string, err ...error) ErrorResponse {
	return NewError(fiber.StatusNotFound, ErrNotFound, message, nil, err...)
}

func EndpointNotFound() ErrorResponse {
	return NewError(fiber.StatusNotFound, ErrNotFound, "Invalid endpoint", nil, nil)
}
