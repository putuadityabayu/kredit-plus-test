/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package test

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"log"
	"os"
	"testing"
	mock_repository "xyz/mocks/repository"
	"xyz/pkg/otel"
)

func setupTestConfig() {
	configFile := `{
  "app_env": "test",
  "secret": {
    "jwt": "61be065c6672292aeb685065264c6c23",
    "password_salt": "password_salt"
  }
}`

	viper.SetConfigType("json")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(configFile)))
	if err != nil {
		log.Fatal(err)
	}

	otel.InitTelemetry(context.Background(), "xyz-test")
}

func teardownTestConfig() {
	viper.Reset()
	otel.Shutdown()
}

func TestMain(m *testing.M) {
	// Setup
	setupTestConfig()

	// Run tests
	code := m.Run()

	teardownTestConfig()

	os.Exit(code)
}

type setupResponse struct {
	ctrl *gomock.Controller

	userRepo  *mock_repository.MockUserRepository
	limitRepo *mock_repository.MockTenorLimitsRepository
}

func setupApp(t *testing.T) *setupResponse {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_repository.NewMockUserRepository(ctrl)
	limitRepo := mock_repository.NewMockTenorLimitsRepository(ctrl)

	userRepo.EXPECT().StartTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		// Jalankan fungsi yang di-pass
		err := fn(ctx)
		return err
	}).AnyTimes()

	return &setupResponse{
		ctrl:      ctrl,
		userRepo:  userRepo,
		limitRepo: limitRepo,
	}
}
