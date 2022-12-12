// Package utils is a package that provides some helper functions.
package utils

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/gophkeeper/internal/server/config"
	"github.com/sergalkin/gophkeeper/pkg/migrations"
)

// CreatePostgresTestConn - creates *pgx.con with testing database.
//
// DSN to testing database is gotten from config.
func CreatePostgresTestConn() *pgx.Conn {
	ctx := context.Background()
	cfg := config.NewConfig()

	con, err := pgx.Connect(ctx, cfg.DSNTest)
	if err != nil {
		panic(err)
	}

	errPing := con.Ping(ctx)
	if errPing != nil {
		panic(errPing)
	}

	return con
}

// RefreshTestDatabase - drops all tables in testing db and then runs all migrations again.
func RefreshTestDatabase() {
	migrator := migrations.NewMigrationManager(config.NewConfig())
	if err := migrator.RefreshTest(); err != nil {
		panic(err)
	}
}
