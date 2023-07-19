package db

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
	"runtime"

	"github.com/ansel1/merry"
	"github.com/doug-martin/goqu/v9"
	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func New() (*goqu.Database, error) {
	if err := doMigration(); err != nil {
		return nil, merry.Wrap(err)
	}

	db, err := sql.Open("postgres", config.DBdsn())
	if err != nil {
		return nil, merry.Wrap(err)
	}

	if err := db.Ping(); err != nil {
		return nil, merry.Wrap(err)
	}

	return goqu.New("postgres", db), nil
}

func doMigration() error {
	if !config.Get.DB.Migrate {
		return nil
	}

	db, err := sql.Open("postgres", config.DBdsn())
	if err != nil {
		return merry.Wrap(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return merry.Wrap(err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		migrationsFile(),
		"postgres",
		driver,
	)
	if err != nil {
		return merry.Wrap(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return merry.Wrap(err)
	}

	return nil
}

func migrationsFile() string {
	_, filename, _, _ := runtime.Caller(1)
	migrationsURL := path.Join(path.Dir(filename), "migrations")
	return fmt.Sprintf("file://%s", migrationsURL)
}

func NewMockDB() (*goqu.Database, error) {
	if err := resetDB(); err != nil {
		return nil, merry.Wrap(err)
	}

	if err := doMigration(); err != nil {
		return nil, merry.Wrap(err)
	}

	db, err := sql.Open("postgres", config.DBdsn())
	if err != nil {
		return nil, merry.Wrap(err)
	}

	return goqu.New("postgres", db), nil
}

func resetDB() error {
	db, err := sql.Open("postgres", config.DBdsn())
	if err != nil {
		return merry.Wrap(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return merry.Wrap(err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		migrationsFile(),
		"postgres",
		driver,
	)
	if err != nil {
		return merry.Wrap(err)
	}

	if err := m.Drop(); err != nil {
		return merry.Wrap(err)
	}

	return nil
}
