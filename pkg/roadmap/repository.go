//go:generate mockery -all -dir ./ -case snake -output ./mocks

package roadmap

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/herr"
	"github.com/peteraba/roadmapper/pkg/repository"
)

type DbReadWriter interface {
	Get(c code.Code) (*Roadmap, error)
	Create(roadmap Roadmap) error
}

// Repository represents a persistence layer using a database (Postgres)
type Repository struct {
	repository.PgRepository
}

// Get retrieves a Roadmap from the database
func (drw Repository) Get(code code.Code) (*Roadmap, error) {
	db := drw.Connect()
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

// Create writes a roadmap to the database
func (drw Repository) Create(roadmap Roadmap) error {
	db := drw.Connect()
	defer db.Close()

	_, err := db.Model(&roadmap).Insert()
	if err != nil {
		return herr.NewFromError(err, http.StatusInternalServerError)
	}

	return err
}
