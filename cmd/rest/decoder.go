/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package rest

import (
	nullable2 "go.portalnesia.com/nullable"
	"reflect"
	"strconv"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/gofiber/fiber/v2"
	"go.portalnesia.com/utils"
)

func registerDecoder() []fiber.ParserType {

	nullBool := fiber.ParserType{
		Customtype: nullable2.Bool{},
		Converter:  boolConverter,
	}
	nullTime := fiber.ParserType{
		Customtype: nullable2.Time{},
		Converter:  timeConverter,
	}
	nullString := fiber.ParserType{
		Customtype: nullable2.String{},
		Converter:  stringConverter,
	}
	nullFloat := fiber.ParserType{
		Customtype: nullable2.Float{},
		Converter:  floatConverter,
	}
	nullInt := fiber.ParserType{
		Customtype: nullable2.Int{},
		Converter:  intConverter,
	}

	return []fiber.ParserType{
		nullBool,
		nullTime,
		nullString,
		nullFloat,
		nullInt,
		nullInt,
	}
}

var timeConverter = func(value string) reflect.Value {
	c := carbon.Parse(value)
	if c.IsValid() {
		a := nullable2.NewTime(c.StdTime(), true, true)
		return reflect.ValueOf(a)
	} else {
		a := nullable2.NewTime(time.Now(), true, false)
		return reflect.ValueOf(a)
	}
}

var boolConverter = func(value string) reflect.Value {
	b := utils.IsTrue(value)
	a := nullable2.NewBool(b, true, true)
	return reflect.ValueOf(a)
}

var stringConverter = func(value string) reflect.Value {
	a := nullable2.NewString(value, true, true)
	return reflect.ValueOf(a)
}

var floatConverter = func(value string) reflect.Value {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		a := nullable2.NewFloat(f, true, false)
		return reflect.ValueOf(a)
	}
	a := nullable2.NewFloat(f, true, true)
	return reflect.ValueOf(a)
}

var intConverter = func(value string) reflect.Value {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		a := nullable2.NewInt(i, true, false)
		return reflect.ValueOf(a)
	}
	a := nullable2.NewInt(i, true, true)
	return reflect.ValueOf(a)
}
