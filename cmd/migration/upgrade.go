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
	upgradeForce  bool = false
	upgradeDryRun bool = false
)

// migrationCmd represents the migration command
var migrationUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade sekala",
	Long:  `Migrate the DB to the most recent version available`,
	Run: func(cmd *cobra.Command, args []string) {
		command := "up"
		if upgradeDryRun {
			command = "validate"
		} else if upgradeForce {
			command = "redo"
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
	migrationCmd.AddCommand(migrationUpgradeCmd)

	migrationUpgradeCmd.Flags().BoolVarP(&upgradeForce, "force", "f", false, "Force re-run the latest migration. DO WITH YOUR OWN RISK")
	migrationUpgradeCmd.Flags().BoolVar(&upgradeDryRun, "dry-run", false, "Check migration files without running them")
}
