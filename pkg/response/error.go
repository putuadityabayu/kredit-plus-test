/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"runtime"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Debug struct {
	TraceID   string `json:"trace_id"`
	Err       error  `json:"error,omitempty"`
	ErrString string `json:"error_string,omitempty"`
}
type ErrorResponse struct {
	Status  string       `json:"status"`
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details"`
	Debug   Debug        `json:"debug"`

	HttpStatus int `json:"-"`
	stack      []uintptr
	frames     []stackFrame
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e ErrorResponse) Response(c *fiber.Ctx) error {
	if appEnv := viper.GetString("app_env"); appEnv == "production" {
		// remove golang err on production
		e.Debug.Err = nil
	} else if e.Debug.Err != nil {
		e.Debug.ErrString = e.Debug.Err.Error()
	}

	return c.Status(e.HttpStatus).JSON(e)
}

func NewError(httpStatus int, code, message string, details []FieldError, err ...error) ErrorResponse {
	stack := make([]uintptr, 5)
	length := runtime.Callers(2, stack)

	debug := Debug{}
	if len(err) > 0 {
		debug.Err = err[0]
	}

	return ErrorResponse{
		Status:     "error",
		Code:       code,
		Message:    message,
		Details:    details,
		Debug:      debug,
		HttpStatus: httpStatus,
		stack:      stack[:length],
	}
}
