package db

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
	"runtime"

	"github.com/doug-martin/goqu/v9"
	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DB *goqu.Database

func New() error {
	if err := doMigration(); err != nil {
		return err
	}

	d, err := sql.Open("postgres", config.DBdsn())
	if err != nil {
		return err
	}

	if err := d.Ping(); err != nil {
		return err
	}

	DB = goqu.New("postgres", d)

	return nil
}

func doMigration() error {
	if !config.Get.DB.Migrate {
		return nil
	}

	db, err := sql.Open("postgres", config.DBdsn())
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		migrationsFile(),
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func migrationsFile() string {
	_, filename, _, _ := runtime.Caller(1)
	migrationsURL := path.Join(path.Dir(filename), "migrations")
	return fmt.Sprintf("file://%s", migrationsURL)
}
