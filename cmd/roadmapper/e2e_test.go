// +build e2e

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/peteraba/roadmapper/pkg/repository"

	"github.com/chromedp/chromedp"
	_ "github.com/lib/pq"
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/roadmap"
	"github.com/peteraba/roadmapper/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	e2eAppPort      uint = 9876
	e2eDbName            = "rdmp"
	e2eDbUser            = "rdmp"
	e2eDbPass            = "secret"
	e2eBaseUrl           = "http://localhost:9876/"
	e2eMatomoDomain      = "https://example.com/matomo"
	e2eDocBaseURL        = "https://docs.rdmp.app/"
	e2eTitle             = "Example Roadmap"
	e2eTxt               = `Monocle ipsum dolor sit amet
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
	e2eBaseURL = "https://example.com/foo"
)

func TestApp_Commandline(t *testing.T) {
	var (
		dateFormat        = "2006-01-02"
		fw, lh     uint64 = 800, 30
		rw                = roadmap.NewIO()
	)

	type args struct {
		rw                  roadmap.IO
		content, output     string
		format              string
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
				e2eTxt,
				"test.svg",
				"svg",
				dateFormat,
				e2eBaseURL,
				fw,
				lh,
			},
		},
		{
			"png size",
			args{
				rw,
				e2eTxt,
				"test.png",
				"png",
				dateFormat,
				e2eBaseURL,
				fw,
				lh,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zap.NewNop()

			err := Render(
				rw,
				logger,
				tt.args.content,
				tt.args.output,
				tt.args.format,
				tt.args.dateFormat,
				tt.args.baseUrl,
				tt.args.fw,
				tt.args.lh,
			)

			require.NoError(t, err)

			expectedData, err := ioutil.ReadFile(fmt.Sprintf("../../res/golden_files/%s", tt.args.output))
			require.NoError(t, err)
			actualData, err := ioutil.ReadFile(tt.args.output)
			require.NoError(t, err)

			ed0, ad0 := float64(len(expectedData)), float64(len(actualData))
			ed1, ad1 := ed0*1.2, ad0*1.2

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

var timeout time.Duration

func getTimeout(t *testing.T) time.Duration {
	if timeout != 0 {
		return timeout
	}

	timeout = 15 * time.Second
	timeoutEnv := os.Getenv("TIMEOUT")
	if timeoutEnv != "" {
		timeoutParsed, err := strconv.ParseInt(timeoutEnv, 10, 32)
		if err != nil {
			t.Errorf("failed parsing TIMEOUT environment variable '%s': %w", timeoutEnv, err)
		}
		timeout = time.Duration(timeoutParsed) * time.Second
	}

	return timeout
}

func TestE2E_Server(t *testing.T) {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, getTimeout(t))
	defer cancel()

	// create a new database
	logger := zap.NewNop()
	baseRepo, teardown := testutils.SetupRepository(t, "TestIntegration_Repository_Get", e2eDbUser, e2eDbPass, e2eDbName, logger)
	defer teardown()

	// start up a new app
	quit := setupApp(t, baseRepo)
	defer teardownApp(quit)

	tests := []struct {
		name     string
		txt      string
		title    string
		baseURL  string
		svgMatch string
		want     string
	}{
		{
			name:     "all filled",
			txt:      e2eTxt,
			title:    e2eTitle,
			baseURL:  e2eBaseURL,
			svgMatch: "Monocle ipsum dolor sit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var txtFound, baseUrlFound, titleFound, svgFound string
			_, _, _ = txtFound, baseUrlFound, svgFound

			err := chromedp.Run(ctx,
				chromedp.Navigate(e2eBaseUrl),
				// wait for form element to become visible (ie, page is loaded)
				chromedp.WaitVisible(`#roadmap-form`),
				// set the value of the title
				chromedp.SetValue(`#title`, tt.title),
				// set the value of the textarea
				chromedp.SetValue(`#txt`, tt.txt),
				// set the value of the base url
				chromedp.SetValue(`#base-url`, tt.baseURL),
				// set the value of the time spent hidden field
				chromedp.SetValue(`#ts`, "10"),
				// submit the form
				chromedp.Submit(`#form-submit`),
				// wait for redirect
				chromedp.WaitVisible(`#roadmap-svg`),
				// retrieve relevant values
				chromedp.Value(`#txt`, &txtFound),
				chromedp.Value(`#base-url`, &baseUrlFound),
				chromedp.Value(`#title`, &titleFound),
				chromedp.WaitVisible(`#roadmap-svg svg`),
				chromedp.OuterHTML(`#roadmap-svg svg`, &svgFound),
			)

			if err != nil {
				t.Fatalf("chromedp run: error = %v", err)
			}

			assert.Equal(t, tt.txt, txtFound)
			assert.Equal(t, tt.baseURL, baseUrlFound)
			assert.Equal(t, tt.title, titleFound)
			if tt.svgMatch != "" {
				assert.Contains(t, svgFound, tt.svgMatch)
			}
		})
	}
}

func setupApp(t *testing.T, baseRepo repository.PgRepository) chan os.Signal {
	quit := make(chan os.Signal, 1)

	_ = t

	rw := roadmap.Repository{PgRepository: baseRepo}
	cb := code.Builder{}

	h := roadmap.NewHandler(baseRepo.Logger, rw, cb, AppVersion, e2eMatomoDomain, e2eDocBaseURL, false)

	go Serve(quit, e2eAppPort, "", "", "../../", h)

	return quit
}

func teardownApp(quit chan os.Signal) {
	quit <- os.Interrupt
}
