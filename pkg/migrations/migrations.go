package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/sergalkin/gophkeeper/internal/server/config"
)

type migrationManager struct {
	cfg config.Config
}

// NewMigrationManager - creates new migration manager instance.
func NewMigrationManager(c config.Config) *migrationManager {
	return &migrationManager{cfg: c}
}

// Up - applies migrations to database via DSN from config.
func (m *migrationManager) Up() error {
	migration, err := migrate.New("file://internal/server/migrations", m.cfg.DSN)
	if err != nil {
		return fmt.Errorf("mingrate.New error: %w", err)
	}
	defer migration.Close()

	migration.Up()

	return nil
}
