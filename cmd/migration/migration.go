/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package migration_cmd

import (
	"xyz/pkg/config"

	"github.com/spf13/cobra"
)

var appConfig config.Config

// migrationCmd represents the migration command
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Database migration",
	Long:  `A tool for migrating databases when upgrading or downgrading Sekala systems`,
}

func Init(appConfigParam config.Config) *cobra.Command {
	appConfig = appConfigParam
	return migrationCmd
}
