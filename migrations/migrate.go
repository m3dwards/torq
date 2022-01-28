package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"net/http"
)

//go:embed *.psql
var static embed.FS

// newMigrationInstance fetches sql files and creates a new migration instance.
func newMigrationInstance(db *sql.DB) (*migrate.Migrate, error) {
	sourceInstance, err := httpfs.New(http.FS(static), ".")
	if err != nil {
		return nil, fmt.Errorf("invalid source instance, %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithInstance("httpfs", sourceInstance, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("could not create migration instance: %v", err)
	}

	return m, nil
}

// MigrateUp migrates up to the latest migration version. It should be used when the version number changes.
func MigrateUp(db *sql.DB) error {
	m, err := newMigrationInstance(db)
	if err != nil {
		return err
	}

	err = m.Up()
	dirtyErr, ok := err.(migrate.ErrDirty)
	// If the Error did not originate from a dirty state, return the error directly.
	if err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion && err != migrate.ErrLocked && !ok {
		return err
	}

	// If the error is due to dirty state. Roll back and try again.
	if ok {
		fmt.Printf("Migration is dirty, forcing rollback and retrying")
		err = m.Force(dirtyErr.Version - 1)
		if err != nil {
			return err
		}
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion && err != migrate.ErrLocked {
			return err
		}
	}

	return nil
}

// MigrateDown migrates the database down one step. Should only be used during development.
func MigrateDown(db *sql.DB) error {
	m, err := newMigrationInstance(db)
	if err != nil {
		return err
	}

	err = m.Steps(-1)
	if err != nil {
		return err
	}

	return nil
}
