package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-pg/pg"
)

// CreateDbReadWriter creates a DbReadWriter instance
func CreateDbReadWriter(applicationName, dbHost, dbPort, dbName, dbUser, dbPass string, logQueries bool) DbReadWriter {
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

	return DbReadWriter{pgOptions: pgOptions, logQueries: logQueries}
}

// DbReadWriter represents a persistence layer using a database (Postgres)
type DbReadWriter struct {
	pgOptions  *pg.Options
	logQueries bool
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
func (dl DbReadWriter) connect() *pg.DB {
	db := pg.Connect(dl.pgOptions)

	if dl.logQueries {
		db.AddQueryHook(dbLogger{})
	}

	return db
}

// Read reads a Roadmap from the database
func (d DbReadWriter) Read(code Code) (*Roadmap, error) {
	db := d.connect()
	defer db.Close()

	r := &Roadmap{ID: code.ID()}

	err := db.Select(r)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, HttpError{error: err, status: http.StatusNotFound}
		}

		return nil, HttpError{error: err, status: http.StatusInternalServerError}
	}

	r.UpdatedAt = time.Now()
	_, err = db.Exec("UPDATE roadmaps SET accessed_at = NOW() WHERE id = ?", r.ID)
	if err != nil {
		return nil, HttpError{error: err, status: http.StatusInternalServerError}
	}

	return r, nil
}

// Write writes a roadmap to the database
func (d DbReadWriter) Write(cb CodeBuilder, roadmap Roadmap) error {
	db := d.connect()
	defer db.Close()

	_, err := db.Model(&roadmap).Insert()
	if err != nil {
		return HttpError{error: err, status: http.StatusInternalServerError}
	}

	return err
}

// CreateFileReadWriter creates a FileReadWriter instance
func CreateFileReadWriter() FileReadWriter {
	return FileReadWriter{}
}

// FileReadWriter represents a persistence layer using the file system (or standard i/o)
type FileReadWriter struct {
}

// Read reads a Roadmap from the file system (or standard i/o)
func (f FileReadWriter) Read(input string) ([]string, error) {
	var (
		file = os.Stdin
		err  error
	)

	if input != "" {
		file, err = os.Open(input)
		if err != nil {
			return nil, fmt.Errorf("can't open file (%s): %w", input, err)
		}
	}

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file (%s): %w", input, err)
	}

	return lines, nil
}

// Write writes a roadmap to the file system (or standard i/o)
func (f FileReadWriter) Write(output string, content string) error {
	if output == "" {
		_, err := fmt.Print(content)

		return err
	}

	d1 := []byte(content)
	err := ioutil.WriteFile(output, d1, 0644)

	return err
}
