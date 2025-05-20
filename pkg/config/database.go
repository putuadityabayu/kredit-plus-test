/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package config

import (
	"fmt"
	"github.com/dromara/carbon/v2"
	othermysql "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	nativeLog "log"
	"os"
	"time"
)

func InitDatabase() *gorm.DB {
	isProduction := viper.GetString("app_env") == "production"
	level := logger.Error
	if isProduction {
		level = logger.Silent
	}

	mysqlconfig := othermysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		DBName:               os.Getenv("DB_NAME"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		ParseTime:            true,
		Loc:                  time.UTC,
		AllowNativePasswords: true,
	}
	dsn := mysqlconfig.FormatDSN()

	portalnesia := mysql.New(mysql.Config{
		DSN: dsn,
	})

	database, err := gorm.Open(portalnesia, &gorm.Config{
		PrepareStmt: true,
		Logger: logger.New(
			nativeLog.New(os.Stdout, "\r\n", nativeLog.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  level,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
		NowFunc: func() time.Time {
			return carbon.Now().StdTime()
		},
	})

	if err != nil {
		log.Fatal("Failed to initialize mysql")
	}

	db, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get mysql.DB")
	}
	db.SetMaxOpenConns(100)                // Jumlah maksimum koneksi yang bisa dibuka
	db.SetMaxIdleConns(10)                 // Jumlah koneksi idle yang dipertahankan
	db.SetConnMaxLifetime(5 * time.Minute) // Masa pakai maksimum koneksi

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping mysql")
	}

	return database
}
