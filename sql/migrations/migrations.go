package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var (
	//go:embed *.sql
	migrationFS embed.FS
)

func Run(db *sql.DB, database string) error {
	sourceInstance, err := iofs.New(migrationFS, ".")
	if err != nil {
		return fmt.Errorf("iofs: %w", err)
	}

	databaseInstance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("sqlite3: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceInstance, database, databaseInstance)
	if err != nil {
		return fmt.Errorf("migrate new: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
