/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package main

import (
	"embed"
	"xyz/cmd"
)

var (
	//go:embed migrations/*
	migrationEmbed embed.FS
)

func main() {
	cmd.Execute(migrationEmbed)
}
