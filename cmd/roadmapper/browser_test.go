// +build e2e

package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/chromedp/chromedp"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/repository"
	"github.com/peteraba/roadmapper/pkg/roadmap"
	"github.com/peteraba/roadmapper/pkg/testutils"
)

func TestE2E_Server(t *testing.T) {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, testutils.GetTimeout(t))
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

	go Serve(quit, e2eAppPort, "", "", "../../static", h)

	return quit
}

func teardownApp(quit chan os.Signal) {
	quit <- os.Interrupt
}
