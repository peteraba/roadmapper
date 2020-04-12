// +build integration

package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/chromedp/chromedp"
	"github.com/labstack/gommon/log"

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

	_, err := migrateUp(dbUser, dbPass, dbHost, dbPort, dbName, 0)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return quit
}

var (
	txt = `Monocle ipsum dolor sit amet
Ettinger punctual izakaya concierge [2020-02-02, 2020-02-20, 60%]
	Zürich Baggu bureaux [/issues/1]
		Toto Comme des Garçons liveable [2020-02-04, 2020-02-25, 100%, /issues/2]
		Winkreative boutique St Moritz [2020-02-06, 2020-02-22, 55%, /issues/3]
	Toto joy perfect Porter [2020-02-25, 2020-03-01, 100%, |1]
Craftsmanship artisanal
	Marylebone exclusive [2020-03-03, 2020-03-10, 100%]
	Beams elegant destination [2020-03-08, 2020-03-12, 100%, |1]
	Winkreative ryokan hand-crafted [2020-03-13, 2020-03-31, 20%]
Nordic Toto first-class Singap
	Concierge cutting-edge Zürich global bureaux
		Sunspel sophisticated lovely uniforms [2020-03-17, 2020-03-31]
		Share blog post on social media [2020-03-17, 2020-03-31, 80%]
	Talk about the tool in relevant meetups [2020-04-01, 2020-06-15, 20%]
Melbourne handsome boutique
	Boutique magna iconic
		Carefully curated laborum destination [2020-03-28, 2020-05-01, 60%]
	Qui incididunt sleepy
		Scandinavian occaecat culpa [2020-03-26, 2020-04-01, 90%]
Hand-crafted K-pop boulevard
	Charming sed quality [2020-03-18, 2020-05-31, 20%]
	Sunspel alluring ut dolore [2020-04-15, 2020-04-30, 30%]
Business class Shinkansen [2020-04-01, 2020-05-31, 45%]
	Nisi excepteur hand-crafted hub
	Ettinger Airbus A380
Essential conversation bespoke
Muji enim

|Laboris ullamco
|Muji enim finest [2020-02-12, https://example.com/abc, bcdef]`
	txtBaseUrl = "https://example.com/foo"
)

func TestApp_TextToRoadmap(t *testing.T) {
	now := time.Now()
	content := Content(txt)

	roadmap := content.ToRoadmap(123, nil, "2006-01-02", txtBaseUrl, now)

	actual := roadmap.String()

	assert.Equal(t, txt, actual)
}

func TestApp_TextToVisual(t *testing.T) {
	now := time.Now()
	content := Content(txt)
	expectedProjectLength := 28
	expectedMilestoneLength := 2
	expectedDeadline1 := time.Date(2020, 3, 12, 0, 0, 0, 0, time.UTC)

	roadmap := content.ToRoadmap(123, nil, "2006-01-02", txtBaseUrl, now)
	visualRoadmap := roadmap.ToVisual()

	assert.Len(t, visualRoadmap.Projects, expectedProjectLength)
	assert.Len(t, visualRoadmap.Milestones, expectedMilestoneLength)
	assert.Equal(t, expectedDeadline1.String(), visualRoadmap.Milestones[0].DeadlineAt.String())
}

func TestApp_Commandline(t *testing.T) {
	var (
		dateFormat        = "2006-01-02"
		fw, lh     uint64 = 800, 30
		rw                = CreateFileReadWriter()
	)

	type args struct {
		rw                  FileReadWriter
		content, output     string
		format              fileFormat
		dateFormat, baseUrl string
		fw, lh              uint64
	}

	tests := []struct {
		name string
		args args
	}{
		{
			"svg size",
			args{
				rw,
				txt,
				"test.svg",
				svgFormat,
				dateFormat,
				txtBaseUrl,
				fw,
				lh,
			},
		},
		{
			"pdf size",
			args{
				rw,
				txt,
				"test.pdf",
				pdfFormat,
				dateFormat,
				txtBaseUrl,
				fw,
				lh,
			},
		},
		{
			"png size",
			args{
				rw,
				txt,
				"test.png",
				pngFormat,
				dateFormat,
				txtBaseUrl,
				fw,
				lh,
			},
		},
		{
			"gif size",
			args{
				rw,
				txt,
				"test.gif",
				gifFormat,
				dateFormat,
				txtBaseUrl,
				fw,
				lh,
			},
		},
		{
			"jpg size",
			args{
				rw,
				txt,
				"test.jpg",
				jpgFormat,
				dateFormat,
				txtBaseUrl,
				fw,
				lh,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Render(
				rw,
				tt.args.content,
				tt.args.output,
				tt.args.format,
				tt.args.dateFormat,
				tt.args.baseUrl,
				tt.args.fw,
				tt.args.lh,
			)

			require.NoError(t, err)

			expectedData, err := ioutil.ReadFile(fmt.Sprintf("goldenfiles/%s", tt.args.output))
			require.NoError(t, err)
			actualData, err := ioutil.ReadFile(tt.args.output)
			require.NoError(t, err)

			ed0, ad0 := float64(len(expectedData)), float64(len(actualData))
			ed1, ad1 := ed0*1.1, ad0*1.1

			assert.Greater(t, ed1, ad0, "generated and golden files differ a lot")
			assert.Less(t, ed0, ad1, "generated and golden files differ a lot")

			if !t.Failed() {
				err = os.Remove(tt.args.output) // remove a single file
				if err != nil {
					t.Errorf("failed to delete file: %s", tt.args.output)
				}
			}
		})
	}
}

func TestApp_Server(t *testing.T) {
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
			name:       "all filled",
			txt:        txt,
			txtBaseUrl: txtBaseUrl,
			svgMatch:   "Monocle ipsum dolor sit",
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
				log.Printf(svgFound)
				assert.Contains(t, svgFound, tt.svgMatch)
			}
		})
	}
}
