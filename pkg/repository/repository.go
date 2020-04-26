package repository

import (
	"fmt"

	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PgRepository struct {
	PgOptions *pg.Options
	Logger    *zap.Logger
}

// NewPgRepository creates a Repository instance
func NewPgRepository(applicationName, dbHost, dbPort, dbName, dbUser, dbPass string, logger *zap.Logger) PgRepository {
	pgOptions := &pg.Options{
		Addr:                  fmt.Sprintf("%s:%s", dbHost, dbPort),
		User:                  dbUser,
		Password:              dbPass,
		Database:              dbName,
		ApplicationName:       applicationName,
		TLSConfig:             nil,
		MaxRetries:            5,
		RetryStatementTimeout: false,
	}

	return PgRepository{PgOptions: pgOptions, Logger: logger}
}

func (pr *PgRepository) InTx(operation func(tx *pg.Tx) error) error {
	db := pr.Connect()
	defer db.Close()

	return db.RunInTransaction(operation)
}

// Connect connects to a database and returns a DB resource
// it can optionally set up the previously query hooks if needed
func (pr PgRepository) Connect() *pg.DB {
	db := pg.Connect(pr.PgOptions)

	if pr.Logger != nil {
		db.AddQueryHook(dbLogger{pr.Logger})
	}

	return db
}

// Connect connects to a database and returns a DB resource
func (pr PgRepository) ConnectNoHook() *pg.DB {
	db := pg.Connect(pr.PgOptions)

	return db
}

type dbLogger struct {
	logger *zap.Logger
}

// BeforeQuery is a go-pg hook that is called before a query is executed
func (dl dbLogger) BeforeQuery(q *pg.QueryEvent) {
	_ = q
}

// BeforeQuery is a go-pg hook that is called after a query is executed
// It's used for logging queries
func (dl dbLogger) AfterQuery(q *pg.QueryEvent) {
	formattedQuery, err := q.FormattedQuery()
	if err != nil {
		dl.logger.Warn("database query error", zap.Error(err))
	} else {
		dl.logger.Info("database query",
			zap.String("formattedQuery", formattedQuery),
		)
	}
}
