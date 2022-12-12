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

// RefreshTest - run DropTest and UpTest functions.
func (m *migrationManager) RefreshTest() error {
	if err := m.DropTest(); err != nil {
		return err
	}
	if err := m.UpTest(); err != nil {
		return err
	}
	return nil
}

// DropTest - drops all tables in testing database.
func (m *migrationManager) DropTest() error {
	migration, err := migrate.New("file://../../../../internal/server/migrations", m.cfg.DSNTest)
	if err != nil {
		return fmt.Errorf("mingrate.New error: %w", err)
	}
	defer migration.Close()

	if err = migration.Drop(); err != nil {
		return fmt.Errorf("error in migrating db: %w", err)
	}

	return nil
}

// UpTest - runs migrations against testing database.
func (m *migrationManager) UpTest() error {
	migration, err := migrate.New("file://../../../../internal/server/migrations", m.cfg.DSNTest)
	if err != nil {
		return fmt.Errorf("mingrate.New error: %w", err)
	}
	defer migration.Close()

	if err = migration.Up(); err != nil {
		return fmt.Errorf("error in migrating db: %w", err)
	}

	return nil
}
