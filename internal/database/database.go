package database

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

//go:embed migrations/*.*.sql
var migrations embed.FS

func OpenDatabase(dbPath string) (*Database, error) {
	if dbPath == "" {
		return nil, errors.New("empty database path")
	}

	sourceDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	databaseDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		db.Close()
		return nil, err
	}

	migration, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", databaseDriver)
	if err != nil {
		db.Close()
		return nil, err
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		return nil, err
	}

	return &Database{db}, nil
}

func CloseDatabase(db *Database) {
	db.db.Close()
}
