// +build integration

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
)

const (
	appPort      uint = 9876
	dbHost            = "localhost"
	dbName            = "rdmp"
	dbUser            = "rdmp"
	dbPass            = "secret"
	baseUrl           = "http://localhost:9876/"
	matomoDomain      = "https://example.com/matomo"
	docBaseUrl        = "https://docs.rdmp.app/"
)

func setupDb(t *testing.T) (*dockertest.Pool, *dockertest.Resource) {
	var db *sql.DB
	var err error

	dbPool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	dbResource, err := dbPool.Run("postgres", "alpine", []string{"POSTGRES_USER=" + dbUser, "POSTGRES_PASSWORD=" + dbPass, "POSTGRES_DB=" + dbName})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	if err = dbPool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, dbResource.GetPort("5432/tcp"), dbName))
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	return dbPool, dbResource
}

func teardownDb(t *testing.T, dbPool *dockertest.Pool, dbResource *dockertest.Resource) {
	if err := dbPool.Purge(dbResource); err != nil {
		t.Fatalf("Could not tear down the database: %s", err)
	}
}

func teardownApp(t *testing.T, quit chan os.Signal) {
	quit <- os.Interrupt
}

func setupApp(t *testing.T, dbResource *dockertest.Resource) chan os.Signal {
	dbPort := dbResource.GetPort("5432/tcp")

	quit := make(chan os.Signal, 1)

	cb := NewCodeBuilder()

	rw := CreateDbReadWriter(applicationName, dbHost, dbPort, dbName, dbUser, dbPass, true)
	go Serve(quit, appPort, "", "", rw, cb, matomoDomain, docBaseUrl, false)

	migrateUp(dbUser, dbPass, dbHost, dbPort, dbName, 0)

	return quit
}

func TestApp_Server(t *testing.T) {
	txt := `Initial development [2020-02-12, 2020-02-20]
Bring website online
	Select and purchase domain [2020-02-04, 2020-02-25, 100%, /issues/1]
	Create server infrastructure [2020-02-25, 2020-02-28, 100%]
Command line tool
	Create backend SVG generation [2020-03-03, 2020-03-10, 100%]
	Replace frontend SVG generation with backend [2020-03-08, 2020-03-12, 100%]
	Create documentation page [2020-03-13, 2020-03-31, 20%]
Marketing
	Create Facebook page [2020-03-17, 2020-03-25, 0%]
	Write blog posts [2020-03-17, 2020-03-31, 0%]
	Share blog post on social media [2020-03-17, 2020-03-31, 0%]
	Talk about the tool in relevant meetups [2020-04-01, 2020-06-15, 0]`

	txtBaseUrl := ""

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// create a new database
	dbPool, dbResource := setupDb(t)
	defer teardownDb(t, dbPool, dbResource)

	// start up a new app
	quit := setupApp(t, dbResource)
	defer teardownApp(t, quit)

	tests := []struct {
		name       string
		txt        string
		txtBaseUrl string
		svgMatch   string
		want       string
	}{
		{
			name:       "empty baseUrl",
			txt:        txt,
			txtBaseUrl: "",
			svgMatch:   "",
		},
		{
			name:       "all filled",
			txt:        txt,
			txtBaseUrl: txtBaseUrl,
			svgMatch:   "Initial development",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var txtFound, txtBaseUrlFound, svgFound string
			_, _, _ = txtFound, txtBaseUrl, svgFound

			err := chromedp.Run(ctx,
				chromedp.Navigate(baseUrl),
				// wait for form element to become visible (ie, page is loaded)
				chromedp.WaitVisible(`#roadmap-form`),
				// set the value of the textarea
				chromedp.SetValue(`#txt`, tt.txt),
				// set the value of the base url
				chromedp.SetValue(`#base-url`, tt.txtBaseUrl),
				// set the value of the base url
				chromedp.Submit(`#form-submit`),
				// wait for redirect
				chromedp.WaitVisible(`#roadmap-svg`),
				// retrieve relevant values
				chromedp.Value(`#txt`, &txtFound),
				chromedp.Value(`#base-url`, &txtBaseUrlFound),
				chromedp.OuterHTML(`#roadmap-svg`, &svgFound),
			)

			if err != nil {
				t.Fatalf("chromedp run: error = %v", err)
			}

			assert.Equal(t, tt.txt, txtFound)
			assert.Equal(t, tt.txtBaseUrl, txtBaseUrlFound)
			if tt.svgMatch != "" {
				assert.Contains(t, svgFound, tt.svgMatch)
			}
		})
	}
}
