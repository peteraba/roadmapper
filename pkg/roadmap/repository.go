//go:generate mockery -all -dir ./ -case snake -output ./mocks

package roadmap

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-pg/pg"
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/herr"
)

type pgRepository struct {
	conn *pg.DB
}

func (r *pgRepository) InTx(operation func(tx *pg.Tx) error) error {
	return r.conn.RunInTransaction(operation)
}

type DbReadWriter interface {
	Get(c code.Code) (*Roadmap, error)
	Write(roadmap Roadmap) error
}

// Repository represents a persistence layer using a database (Postgres)
type Repository struct {
	pgRepository
	pgOptions  *pg.Options
	logQueries bool
}

// NewRepository creates a Repository instance
func NewRepository(applicationName, dbHost, dbPort, dbName, dbUser, dbPass string, logQueries bool) Repository {
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

	return Repository{pgOptions: pgOptions, logQueries: logQueries}
}

type dbLogger struct{}

// BeforeQuery is a go-pg hook that is called before a query is executed
func (dl dbLogger) BeforeQuery(q *pg.QueryEvent) {
}

// BeforeQuery is a go-pg hook that is called after a query is executed
// It's used for logging queries
func (dl dbLogger) AfterQuery(q *pg.QueryEvent) {
	formattedQuery, _ := q.FormattedQuery()
	fmt.Println(formattedQuery)
}

// connect ensures connects to a database
// it can optionally set up the previously query hooks if needed
func (drw Repository) connect() *pg.DB {
	db := pg.Connect(drw.pgOptions)

	if drw.logQueries {
		db.AddQueryHook(dbLogger{})
	}

	return db
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

// Write writes a roadmap to the database
func (drw Repository) Write(roadmap Roadmap) error {
	db := drw.connect()
	defer db.Close()

	_, err := db.Model(&roadmap).Insert()
	if err != nil {
		return herr.NewFromError(err, http.StatusInternalServerError)
	}

	return err
}
