//go:generate mockery -all -dir ./ -case snake -output ./mocks

package roadmap

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/herr"
	"go.uber.org/zap"
)

type pgRepository struct {
	pgOptions *pg.Options
	logger    *zap.Logger
}

func (r *pgRepository) InTx(operation func(tx *pg.Tx) error) error {
	db := r.connect()
	defer db.Close()

	return db.RunInTransaction(operation)
}

// connect ensures connects to a database
// it can optionally set up the previously query hooks if needed
func (pr pgRepository) connect() *pg.DB {
	db := pg.Connect(pr.pgOptions)

	if pr.logger != nil {
		db.AddQueryHook(dbLogger{pr.logger})
	}

	return db
}

type dbLogger struct {
	logger *zap.Logger
}

// BeforeQuery is a go-pg hook that is called before a query is executed
func (dl dbLogger) BeforeQuery(q *pg.QueryEvent) {
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

type DbReadWriter interface {
	Get(c code.Code) (*Roadmap, error)
	Upsert(roadmap Roadmap) error
}

// Repository represents a persistence layer using a database (Postgres)
type Repository struct {
	pgRepository
}

// NewRepository creates a Repository instance
func NewRepository(applicationName, dbHost, dbPort, dbName, dbUser, dbPass string, logger *zap.Logger) Repository {
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
	pgRepo := pgRepository{pgOptions: pgOptions, logger: logger}

	return Repository{pgRepository: pgRepo}
}

// Get retrieves a Roadmap from the database
func (drw Repository) Get(code code.Code) (*Roadmap, error) {
	db := drw.connect()
	defer db.Close()

	roadmap := &Roadmap{ID: code.ID()}

	err := db.Select(roadmap)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, herr.NewFromError(err, http.StatusNotFound)
		}

		return nil, herr.NewFromError(err, http.StatusInternalServerError)
	}

	roadmap.UpdatedAt = time.Now()
	_, err = db.Exec("UPDATE roadmaps SET accessed_at = NOW() WHERE id = ?", roadmap.ID)
	if err != nil {
		return nil, herr.NewFromError(err, http.StatusInternalServerError)
	}

	return roadmap, nil
}

// Upsert writes a roadmap to the database
func (drw Repository) Upsert(roadmap Roadmap) error {
	db := drw.connect()
	defer db.Close()

	_, err := db.Model(&roadmap).Insert()
	if err != nil {
		return herr.NewFromError(err, http.StatusInternalServerError)
	}

	return err
}
