package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func getConnStr(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	return connStr
}

func getConn(dbUser, dbPass, dbHost, dbPort, dbName string) (*sql.DB, error) {
	connStr := getConnStr(dbUser, dbPass, dbHost, dbPort, dbName)

	return sql.Open("postgres", connStr)
}

func migrateDown(dbUser, dbPass, dbHost, dbPort, dbName string, steps int) (int, error) {
	db, err := getConn(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		return 0, fmt.Errorf("failed to create a connection: %w", err)
	}

	source := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}

	n, err := migrate.ExecMax(db, "postgres", source, migrate.Down, steps)
	if err != nil {
		return 0, fmt.Errorf("migration failed: %w", err)
	}

	return n, nil
}

func migrateUp(dbUser, dbPass, dbHost, dbPort, dbName string, steps int) (int, error) {
	db, err := getConn(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		return 0, fmt.Errorf("failed to create a connection: %w", err)
	}

	source := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}

	n, err := migrate.ExecMax(db, "postgres", source, migrate.Up, steps)
	if err != nil {
		return 0, fmt.Errorf("migration failed: %w", err)
	}

	return n, nil
}
