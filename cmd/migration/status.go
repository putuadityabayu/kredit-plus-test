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

// migrationCmd represents the migration command
var migrationStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Migration status",
	Long:  `Dump the migration status for the current DB`,
	Run: func(cmd *cobra.Command, args []string) {
		dbPkg := config.InitDatabase()
		db, err := dbPkg.DB()
		if err != nil {
			log.Fatal("Failed to get database connection")
		}
		migration.New(db, appConfig.MigrationEmbed).Run(context.TODO(), "status")
	},
}

func init() {
	migrationCmd.AddCommand(migrationStatusCmd)
}
