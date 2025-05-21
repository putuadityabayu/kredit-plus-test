/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package migration

import (
	"context"
	"database/sql"
	"embed"
	"github.com/gofiber/fiber/v2/log"

	"github.com/pressly/goose/v3"
)

type Goose struct {
	db  *sql.DB
	dir string
}

func New(db *sql.DB, migrationFs embed.FS) *Goose {
	goose.SetBaseFS(migrationFs)
	goose.SetTableName("db_version")
	if err := goose.SetDialect(string(goose.DialectMySQL)); err != nil {
		log.Fatal("Failed to set dialect to mysql")
	}

	return &Goose{db: db, dir: "migrations"}
}

func (g *Goose) Run(ctx context.Context, command string, args ...string) {
	if err := goose.RunContext(ctx, command, g.db, g.dir, args...); err != nil {
		log.Fatalf("Command %s: %s", command, err.Error())
	}
}
