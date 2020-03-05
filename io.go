package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-pg/pg"
)

type ReadWriter interface {
	Read(code Code) ([]string, error)
	Write(cb CodeBuilder, code Code, content string) error
}

func CreateReadWriter(dbHost, dbPort, dbName, dbUser, dbPass string) ReadWriter {
	if dbHost != "" && dbPort != "" && dbName != "" && dbUser != "" {
		pgOptions := &pg.Options{
			Addr:                  fmt.Sprintf("%s:%s", dbHost, dbPort),
			User:                  dbUser,
			Password:              dbPass,
			Database:              dbName,
			ApplicationName:       "roadmapper",
			TLSConfig:             nil,
			MaxRetries:            5,
			RetryStatementTimeout: false,
		}
		return DbReadWriter{pgOptions: pgOptions}
	}

	return FileReadWriter{}
}

type DbReadWriter struct {
	pgOptions *pg.Options
}

type roadmap struct {
	Id     int64
	PrevId int64
	Txt    string
}

func (d DbReadWriter) Read(code Code) ([]string, error) {
	db := pg.Connect(d.pgOptions)
	defer db.Close()

	r := &roadmap{Id: code.ID()}

	err := db.Select(r)
	if err != nil {
		return nil, err
	}

	return strings.Split(r.Txt, "\n"), nil
}

func (d DbReadWriter) Write(cb CodeBuilder, code Code, content string) error {
	db := pg.Connect(d.pgOptions)
	defer db.Close()

	// we must find a code that does not yet exist
	var newCode Code
	var found bool
	for {
		newCode = cb.New()
		_, err := d.Read(newCode)
		if err != nil {
			found = true
			break
		}
	}

	if !found {
		return errors.New("no new code found during insert")
	}

	r := &roadmap{Id: newCode.ID(), PrevId: code.ID(), Txt: content}

	err := db.Insert(r)

	return err
}

type FileReadWriter struct {
}

func (f FileReadWriter) Read(code Code) ([]string, error) {
	file, err := os.Open(code.String())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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

func (f FileReadWriter) Write(cb CodeBuilder, code Code, content string) error {
	d1 := []byte(content)
	err := ioutil.WriteFile(code.String(), d1, 0644)

	return err
}
