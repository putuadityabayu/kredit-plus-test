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
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"regexp"
)

const (
	ErrNotFound = "NOT_FOUND"
)

func NotFound(message string, err ...error) ErrorResponse {
	return NewError(fiber.StatusNotFound, ErrNotFound, message, nil, err...)
}

func EndpointNotFound() ErrorResponse {
	return NewError(fiber.StatusNotFound, ErrNotFound, "Invalid endpoint", nil, nil)
}

type SpanInterface interface {
	RecordErrorHelper(err error, message string)
}

func NotfoundHelper(err error, message string, span ...SpanInterface) ErrorResponse {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(message, err)
	}
	if len(span) > 0 {
		span[0].RecordErrorHelper(err, "NotfoundHelper")
	}
	return ErrorServer(MsgInternalServer, err)
}

func parseDuplicateValue(mysqlErr *mysql.MySQLError) (string, string, bool) {
	// Contoh pesan: Duplicate entry '1234567890123456' for key 'idx_users_nik'
	re := regexp.MustCompile(`Duplicate entry '(.+)' for key '(.+)'`)
	matches := re.FindStringSubmatch(mysqlErr.Message)
	if len(matches) == 3 {
		value := matches[1]
		index := matches[2]
		return value, index, true
	}
	return "", "", false
}

// DatabaseHelper is helper to unique constraint error
func DatabaseHelper(err error, mapIndexKey map[string]string, span ...SpanInterface) ErrorResponse {
	var e *mysql.MySQLError
	if errors.As(err, &e) {
		if e.Number == 1062 {
			val, key, ok := parseDuplicateValue(e)
			if ok {
				if index, ok := mapIndexKey[key]; ok {
					return NewError(fiber.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", fmt.Sprintf("%s '%s' already exists", index, val), nil, err)
				} else {
					return NewError(fiber.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", fmt.Sprintf("%s already exists", val), nil, err)
				}
			}
		}
	}

	if len(span) > 0 {
		span[0].RecordErrorHelper(err, "DatabaseHelper")
	}
	return ErrorServer(MsgInternalServer, err)
}
