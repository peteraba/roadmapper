package migrations

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/peteraba/roadmapper/pkg/bindata"
	migrate "github.com/rubenv/sql-migrate"
)

func getConn(dbUser, dbPass, dbHost, dbPort, dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	return sql.Open("postgres", connStr)
}

func Down(dbUser, dbPass, dbHost, dbPort, dbName string, steps int) (int, error) {
	db, err := getConn(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		return 0, fmt.Errorf("failed to create a connection: %w", err)
	}

	source := &migrate.AssetMigrationSource{
		Asset:    bindata.Asset,
		AssetDir: bindata.AssetDir,
		Dir:      "migrations",
	}

	n, err := migrate.ExecMax(db, "postgres", source, migrate.Down, steps)
	if err != nil {
		return 0, fmt.Errorf("migration failed: %w", err)
	}

	return n, nil
}

func Up(dbUser, dbPass, dbHost, dbPort, dbName string, steps int) (int, error) {
	db, err := getConn(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		return 0, fmt.Errorf("failed to create a connection: %w", err)
	}

	source := &migrate.AssetMigrationSource{
		Asset:    bindata.Asset,
		AssetDir: bindata.AssetDir,
		Dir:      "migrations",
	}

	n, err := migrate.ExecMax(db, "postgres", source, migrate.Up, steps)
	if err != nil {
		return 0, fmt.Errorf("migration failed: %w", err)
	}

	return n, nil
}
