/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package migration_cmd

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"xyz/pkg/config"
	"xyz/pkg/migration"

	"github.com/spf13/cobra"
)

var downgradeAll bool = false

// migrationCmd represents the migration command
var migrationDowngradeCmd = &cobra.Command{
	Use:   "downgrade",
	Short: "Downgrade sekala",
	Long:  `Rollback database migration version`,
	Run: func(cmd *cobra.Command, args []string) {
		command := "down"
		if downgradeAll {
			command = "reset"
		}
		dbPkg := config.InitDatabase()
		db, err := dbPkg.DB()
		if err != nil {
			log.Fatal("Failed to get database connection")
		}

		migration.New(db, appConfig.MigrationEmbed).Run(context.TODO(), command)
	},
}

func init() {
	migrationCmd.AddCommand(migrationDowngradeCmd)

	migrationDowngradeCmd.Flags().BoolVarP(&downgradeAll, "all", "a", false, "Reset all migration. DO WITH YOUR OWN RISK")
}
