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
	"os"
	"testing"
	mock_repository "xyz/mock/repository"
	"xyz/pkg/otel"
)

func setupTestConfig() {
	configFile := `{
  "app_env": "test",
  "secret": {
    "jwt": "61be065c6672292aeb685065264c6c23",
    "password_salt": "password_salt"
  },
}`

	viper.SetConfigType("json")
	_ = viper.ReadConfig(bytes.NewBuffer([]byte(configFile)))

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

	os.Exit(code)
}

type setupResponse struct {
	ctrl *gomock.Controller

	userRepo *mock_repository.MockUserRepository
}

func setupApp(t *testing.T) *setupResponse {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_repository.NewMockUserRepository(ctrl)

	userRepo.EXPECT().StartTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		// Jalankan fungsi yang di-pass
		err := fn(ctx)
		return err
	}).AnyTimes()

	return &setupResponse{
		ctrl:     ctrl,
		userRepo: userRepo,
	}
}
