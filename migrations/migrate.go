package migrations

import (
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"net/http"
)

//go:embed *.psql
var static embed.FS

// newMigrationInstance fetches sql files and creates a new migration instance.
func newMigrationInstance() (*migrate.Migrate, error) {
	sourceInstance, err := httpfs.New(http.FS(static), ".")
	if err != nil {
		return nil, fmt.Errorf("invalid source instance, %w", err)
	}

	m, err := migrate.NewWithSourceInstance(
		"httpfs",
		sourceInstance,
		"postgres://torq:password@localhost:5432/torq?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("Could not create migration instance: %v", err)
	}

	return m, nil
}

// MigrateUp migrates up to the latest migration version. It should be used when the version number changes.
func MigrateUp() error {
	m, err := newMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}

// MigrateDown migrates the database down one step. Should only be used during development.
func MigrateDown() error {
	m, err := newMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Steps(-1)
	if err != nil {
		return err
	}

	return nil
}
