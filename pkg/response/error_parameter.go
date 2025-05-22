/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

const (
	ErrBadRequest     = "BAD_REQUEST"
	MsgInvalidRequest = "Invalid request parameter"
	ErrUnprocessable  = "UNPROCESSABLE_ENTITY"
)

type ErrorFields []FieldError

func (e *ErrorFields) Error() string {
	if e == nil {
		return "empty error"
	}

	buff := bytes.NewBufferString("")
	v := *e
	for i := 0; i < len(v); i++ {
		buff.WriteString(fmt.Sprintf("%s: %s", v[i].Field, v[i].Message))
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func (e *ErrorFields) Add(field string, message string) {
	if e != nil {
		*e = append(*e, FieldError{
			Field:   field,
			Message: message,
		})
	}
}

func (e *ErrorFields) Exist() bool {
	if e == nil {
		return false
	}
	return len(*e) > 0
}

func NewErrorFields(fields ...[2]string) *ErrorFields {
	errFields := make(ErrorFields, len(fields))
	for i, field := range fields {
		errFields[i] = FieldError{
			Field:   field[0],
			Message: field[1],
		}
	}
	return &errFields
}

// ErrorParameter
//
// options:
// - int: http status code
// - error error fields
func ErrorParameter(code string, msg string, options ...any) ErrorResponse {
	var (
		httpStatus = fiber.StatusBadRequest
		fields     []FieldError
		errs       error
	)
	if len(options) > 0 {
		for _, opt := range options {
			switch o := opt.(type) {
			case int:
				httpStatus = o
			case error:

				fieldError := false
				var validateError validator.ValidationErrors
				if errors.As(o, &validateError) {
					fieldError = true
					for _, v := range validateError {
						fields = append(fields, FieldError{
							Field:   v.Field(),
							Message: msgForTag(v),
						})
					}
				}
				var errorField *ErrorFields
				if errors.As(o, &errorField) {
					fieldError = true
					fields = append(fields, *errorField...)
				}
				if !fieldError {
					errs = errors.Join(errs, o)
				}
			}
		}
	}

	return NewError(httpStatus, code, msg, fields, errs)
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
