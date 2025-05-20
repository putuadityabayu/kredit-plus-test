/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package config

import (
	"embed"
	"fmt"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"log"
	"strings"
)

type Config struct {
	MigrationEmbed embed.FS
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	_ = gotenv.Load()

	// Init config
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("xyz")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetDefault("app_env", "local")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Failed to load config file")
	}

	fmt.Println("app_env", viper.AllSettings())
}
