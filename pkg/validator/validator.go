/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package validator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func New() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("field"), ",", 2)[0]
		if name == "" || name == "-" {
			name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		}
		if name == "" || name == "-" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}
		if name == "" || name == "-" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	return validate
}
