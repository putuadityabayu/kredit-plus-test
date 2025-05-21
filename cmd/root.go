/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package cmd

import (
	"embed"
	"log"
	migration_cmd "xyz/cmd/migration"
	"xyz/pkg/config"

	"github.com/spf13/cobra"
)

var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "xyz",
	Short:   "xyz api server",
	Long:    `Aplikasi API Sistem Pembiayaan PT XYZ untuk manajemen konsumen, limit, dan transaksi.`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(migrationEmbed embed.FS) {
	cfg.MigrationEmbed = migrationEmbed

	rootCmd.AddCommand(migration_cmd.Init(cfg))

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)
}
