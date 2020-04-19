package migrations

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/peteraba/roadmapper/pkg/bindata"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrations struct {
	conn *sql.DB
}

func New(dbUser, dbPass, dbHost, dbPort, dbName string) *Migrations {
	conn, err := getConn(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		panic(err)
	}

	m := Migrations{
		conn: conn,
	}

	return &m
}

func NewFromConn(conn *sql.DB) *Migrations {
	m := Migrations{
		conn: conn,
	}

	return &m
}

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

func (m *Migrations) Down(steps int) (int, error) {
	source := &migrate.AssetMigrationSource{
		Asset:    bindata.Asset,
		AssetDir: bindata.AssetDir,
		Dir:      "res/migrations",
	}

	n, err := migrate.ExecMax(m.conn, "postgres", source, migrate.Down, steps)
	if err != nil {
		return 0, fmt.Errorf("migration failed: %w", err)
	}

	return n, nil
}

func (m *Migrations) Up(steps int) (int, error) {
	source := &migrate.AssetMigrationSource{
		Asset:    bindata.Asset,
		AssetDir: bindata.AssetDir,
		Dir:      "res/migrations",
	}

	n, err := migrate.ExecMax(m.conn, "postgres", source, migrate.Up, steps)
	if err != nil {
		return 0, fmt.Errorf("migration failed: %w", err)
	}

	return n, nil
}
