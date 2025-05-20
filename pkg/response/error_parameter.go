/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

const (
	ErrBadRequest        = "BAD_REQUEST"
	ErrInsufficientLimit = "INSUFFICIENT_LIMIT"
)

func ErrorParameter(code string, msg string, err error, status ...int) ErrorResponse {
	var (
		httpStatus = fiber.StatusBadRequest
		fields     []FieldError
	)
	if len(status) > 0 {
		httpStatus = status[0]
	}

	var validateError validator.ValidationErrors
	if errors.As(err, &validateError) {
		for _, v := range validateError {
			fields = append(fields, FieldError{
				Field:   v.Field(),
				Message: msgForTag(v),
			})
		}
	}

	return NewError(httpStatus, code, msg, fields, err)
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Missing parameter `%s`", fe.Field())
	case "email":
		return "Invalid email"
	case "oneof":
		return fmt.Sprintf("Parameter `%s` must be one of: %s", fe.Field(), strings.ReplaceAll(fe.Param(), " ", ", "))
	case "min":
		return fmt.Sprintf("Parameter `%s` must have at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("Parameter `%s` must have a maximum of %s characters", fe.Field(), fe.Param())
	}
	return fmt.Sprintf("Invalid `%s` parameter", fe.Field())
}
