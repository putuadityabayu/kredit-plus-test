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

var (
	createName string
	createType string = "sql"
	createFix  bool   = false
)

// migrationCmd represents the migration command
var migrationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create migration (Dev only)",
	Long: `Creates new migration file with the current timestamp.
Only in development mode`,
	Run: func(cmd *cobra.Command, args []string) {
		command := "create"
		var arg []string
		if createFix {
			command = "fix"
		} else {
			if createName == "" {
				log.Error("Missing file name")
				return
			}
			arg = []string{createName, createType}
		}

		dbPkg := config.InitDatabase()

		db, err := dbPkg.DB()
		if err != nil {
			log.Fatal("Failed to get database connection")
		}

		migration.New(db, appConfig.MigrationEmbed).Run(context.TODO(), command, arg...)
	},
}

func init() {
	migrationCmd.AddCommand(migrationCreateCmd)

	migrationCreateCmd.Flags().StringVar(&createName, "name", "", "Migration name")
	migrationCreateCmd.Flags().StringVar(&createType, "type", "sql", "File type. Valid values are sql or go")
	migrationCreateCmd.Flags().BoolVar(&createFix, "fix", false, "Apply sequential ordering to migrations")
}
