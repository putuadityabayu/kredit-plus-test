/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"runtime/debug"
	"time"
	"xyz/internal/middleware"
	"xyz/internal/repository"
	"xyz/internal/router"
	"xyz/pkg/config"
	"xyz/pkg/otel"
	"xyz/pkg/response"
)

type Rest struct {
	fiberApp *fiber.App
	db       *gorm.DB
}

func New(ctx context.Context) *Rest {
	otel.InitTelemetry(ctx, "xyz-api")
	db := config.InitDatabase()
	fiberStorage := config.InitFiberStorage()

	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType:        registerDecoder(),
		ZeroEmpty:         true,
	})

	app := fiber.New(fiber.Config{
		AppName: "XYZ API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			span := otel.FromContext(c.UserContext())
			defer span.End()

			var e response.ErrorResponse
			if errors.As(err, &e) {
				e.Debug.TraceID = span.GetTraceID()
			} else {
				e = response.ErrorServer(response.MsgInternalServer, err)
				e.Debug.TraceID = span.GetTraceID()
			}
			return e.Response(c)
		},
	})

	app.Use(recover2.New(recover2.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			span := otel.FromContext(c.UserContext())
			defer span.End()

			if err, ok := e.(error); ok {
				span.RecordErrorHelper(err, "error in fiber.StackTraceHandler")
			} else { // not error
				err = errors.New("stack-trace-error")
				span.RecordErrorHelper(err, fmt.Sprintf("%v", e))
			}
			log.Errorf("panic: %v\n%s\n", e, debug.Stack())
		},
	}))

	app.Use(requestid.New())
	app.Use(etag.New())
	app.Use(earlydata.New())
	app.Use(idempotency.New())

	// Inject otel
	app.Use(func(c *fiber.Ctx) error {
		_, span := otel.StartSpanHandler(c, "Request")
		defer span.End()
		return c.Next()
	})

	app.Use(middleware.RateLimit(fiberStorage))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"health": "ok",
		})
	})

	// REPO
	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	repoRegistry := repository.RepoRegistry{
		UserRepository:        userRepo,
		TransactionRepository: transactionRepo,
	}

	// ROUTER
	router.UserRouterV1(app, repoRegistry)
	router.AuthRouterV1(app, repoRegistry)
	router.TransactionRouterV1(app, repoRegistry)

	app.Use(func(c *fiber.Ctx) error {
		return response.EndpointNotFound().Response(c)
	})

	go app.Listen(":" + viper.GetString("port"))

	return &Rest{
		fiberApp: app,
		db:       db,
	}
}

func (r *Rest) Shutdown() {
	otel.Shutdown()
	if err := r.fiberApp.ShutdownWithTimeout(time.Second * 5); err != nil {
		log.Error("Error when shutting down server")
	}
	if r.db != nil {
		if sqlDB, _ := r.db.DB(); sqlDB != nil {
			if err := sqlDB.Close(); err != nil {
				log.Errorf("error closing postgre connection: %s", err.Error())
			}
		}
	}
}
