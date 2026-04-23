package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type (
	// Config holds information about a Database config.
	// Used and formatted for DSN string.
	Config struct {
		Host     string
		Port     string
		Username string
		Password string
		Database string
	}

	// Database holds information about a Database (setup).
	Database struct {
		log    *log.Logger
		config Config
		driver string
	}
)

// New returns new Database.
func New(log *log.Logger, driver string, config Config) *Database {
	return &Database{
		log:    log,
		driver: driver,
		config: config,
	}
}

// SetupDB opens and ping Database using the DSN created by config.
func (d *Database) SetupDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("mysql", d.dsn())
	if err != nil {
		return nil, err
	}

	return db, db.PingContext(ctx)
}

// RunMigrations runs migrations from location.
func (d *Database) RunMigrations(migrationPath string) error {
	absPath, ok, err := validateAbsPath(migrationPath)
	if !ok || err != nil {
		return err
	}

	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", absPath),
		d.dsn(),
	)
	if err != nil {
		return err
	}
	defer func(migrator *migrate.Migrate) {
		if err, _ = migrator.Close(); err != nil {
			d.log.Print(fmt.Sprintf("migrations failed: %v", err))
		}
	}(migrator)

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

//
// Unexported helper functions.
//

func (d *Database) dsn() string {
	sqlCfg := mysql.Config{
		User:   d.config.Username,
		Passwd: d.config.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", d.config.Host, d.config.Port),
		DBName: d.config.Database,
	}

	return sqlCfg.FormatDSN()
}

func validateAbsPath(path string) (string, bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", false, err
	}

	fi, err := os.Stat(absPath)
	if err != nil {
		return "", false, err
	}

	if !fi.IsDir() {
		return "", false, nil
	}

	return absPath, true, nil
}
