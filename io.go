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

type DbReadWriter struct {
	pgOptions  *pg.Options
	logQueries bool
}

type dbLogger struct{}

func (dl dbLogger) BeforeQuery(q *pg.QueryEvent) {
}

func (dl dbLogger) AfterQuery(q *pg.QueryEvent) {
	formattedQuery, _ := q.FormattedQuery()
	fmt.Println(formattedQuery)
}

func (dl DbReadWriter) connect() *pg.DB {
	db := pg.Connect(dl.pgOptions)

	if dl.logQueries {
		db.AddQueryHook(dbLogger{})
	}

	return db
}

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

func (d DbReadWriter) Write(cb CodeBuilder, roadmap Roadmap) error {
	db := d.connect()
	defer db.Close()

	_, err := db.Model(&roadmap).Insert()
	if err != nil {
		return HttpError{error: err, status: http.StatusInternalServerError}
	}

	return err
}

func CreateFileReadWriter() FileReadWriter {
	return FileReadWriter{}
}

type FileReadWriter struct {
}

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

func (f FileReadWriter) Write(output string, content string) error {
	if output == "" {
		_, err := fmt.Print(content)

		return err
	}

	d1 := []byte(content)
	err := ioutil.WriteFile(output, d1, 0644)

	return err
}
