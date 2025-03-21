package sqlite3

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func runMigrations(db *sql.DB) error {
	driver, err := WithInstance(db, &Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}
	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs source: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	if dirty, err := isDirty(m); err != nil {
		return fmt.Errorf("failed to check if database is dirty: %w", err)
	} else if dirty {
		log.Println("Database is dirty, resetting...")
		if err := rerunMigrations(m); err != nil {
			return fmt.Errorf("failed to rerun migrations: %w", err)
		}
	} else {
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}

func isDirty(m *migrate.Migrate) (bool, error) {
	_, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return false, nil
		}
		return false, err
	}
	return dirty, nil
}

func resetDatabase(dbPath string) error {
	if err := os.Remove(dbPath); err != nil {
		return err
	}
	_, err := os.Create(dbPath)
	return err
}

func rerunMigrations(m *migrate.Migrate) error {
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
