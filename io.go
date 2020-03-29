package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

type roadmap struct {
	Id         int64
	PrevId     int64
	Txt        string
	DateFormat string
	BaseUrl    string
	UpdatedAt  time.Time
	AccessedAt time.Time
}

func (d DbReadWriter) Read(code Code) ([]string, string, string, error) {
	db := d.connect()
	defer db.Close()

	r := &roadmap{Id: code.ID()}

	err := db.Select(r)
	if err != nil {
		return nil, "", "", err
	}

	r.UpdatedAt = time.Now()
	_, err = db.Exec("UPDATE roadmaps SET accessed_at = NOW() WHERE id = ?", r.Id)
	if err != nil {
		return nil, "", "", err
	}

	return strings.Split(r.Txt, "\n"), r.DateFormat, r.BaseUrl, nil
}

func (d DbReadWriter) Write(cb CodeBuilder, code Code, content, dateFormat, baseUrl string) (Code, error) {
	db := d.connect()
	defer db.Close()

	// we must find a code that does not yet exist
	newCode := cb.New()
	found := false
	for i := 0; i < 100; i++ {
		_, _, _, err := d.Read(newCode)
		if err != nil {
			found = true
			break
		}
		newCode = cb.New()
	}

	if !found {
		return nil, errors.New("no new code found during insert")
	}

	var prevID int64
	if code != nil {
		prevID = code.ID()
	}

	r := &roadmap{Id: newCode.ID(), PrevId: prevID, Txt: content, DateFormat: dateFormat, BaseUrl: baseUrl, UpdatedAt: time.Now()}

	err := db.Insert(r)

	return newCode, err
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
			return nil, err
		}
	}

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
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
