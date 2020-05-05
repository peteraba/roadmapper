// +build e2e

package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/chromedp/chromedp"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/peteraba/roadmapper/pkg/testutils"
)

const (
	e2eAppPort uint = 9876
	e2eDbName       = "rdmp"
	e2eDbUser       = "rdmp"
	e2eDbPass       = "secret"
	e2eBaseUrl      = "http://localhost:9876/"
	e2eTitle        = "How To Start a Startup"
	e2eTxt          = `Find the idea [2019-07-20, 2020-01-20, 100%]
	Look for things missing in life
	Formalize your idea, run thought experiments [https://example.com/initial-plans]
	Survey friends, potential users or customers [https://example.com/survey-results]
	Go back to the drawing board [https://example.com/reworked-plans]
Validate the idea [2020-01-21, 2020-04-20]
	Make a prototype #1 [2020-01-21, 2020-04-10, 100%, TCK-1, https://github.com/peteraba/roadmapper, |1]
	Show the prototype to 100 people #1 [2020-04-11, 2020-04-20, 80%, TCK-123]
	Analyse results [2020-04-21, 2020-05-05]
	Improve prototype [2020-05-06, 2020-06-06]
	Show the prototype to 100 people #2 [2020-06-07, 2020-06-16]
	Analyse results [2020-06-16, 2020-06-30]
	Improve prototype [2020-07-01, 2020-07-16]
	Show the prototype to 100 people #2 [2020-07-17, 2020-07-25]
Start a business
	Learn about your options about various company types [2019-07-20, 2020-08-31]
	Learn about your options for managing equity [2019-07-20, 2020-08-01]
	Find a co-founder [2020-04-20, 2020-08-31]
	Register your business [2020-08-01, 2020-09-30, |2]
	Look for funding [2020-08-01, 2020-10-31]
	Build a team [2020-11-01, 2020-12-15]
Build version one [2021-01-01, 2021-04-15]
	Build version one [2021-01-01, 2021-03-31]
	Launch [2021-04-01, 2021-04-15, |3]
Grow [2021-04-16, 2021-12-31]
	Follow up with users
	Iterate / Pivot
	Launch again
	Get to 1,000 users
	Plan next steps

|Create the first prototype
|Start your business
|Lunch version one`
	e2eBaseURL = "https://example.com/foo"
)

func TestE2E_Browser(t *testing.T) {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, testutils.GetTimeout(t))
	defer cancel()

	// create a new logger
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// create a new database
	baseRepo, reset, dbTeardown := testutils.SetupRepository(t, "TestE2E_Browser", e2eDbUser, e2eDbPass, e2eDbName, logger)
	defer dbTeardown()

	// start up a new app
	appTeardown := testApp(baseRepo, logger, e2eAppPort, "../../static")
	defer appTeardown()

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
			svgMatch: "Find the idea",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer reset()

			var txtFound, baseUrlFound, bodyFound, titleFound, svgFound string

			_, _, _, _, _ = txtFound, baseUrlFound, bodyFound, titleFound, svgFound

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
				// // retrieve relevant values
				chromedp.Value(`#txt`, &txtFound),
				chromedp.Value(`#base-url`, &baseUrlFound),
				chromedp.Value(`#title`, &titleFound),
				chromedp.WaitVisible(`#roadmap-svg svg`),
				chromedp.OuterHTML(`#roadmap-svg svg`, &svgFound),
			)

			assert.Equal(t, tt.txt, txtFound)
			assert.Equal(t, tt.baseURL, baseUrlFound)
			assert.Equal(t, tt.title, titleFound)
			if tt.svgMatch != "" {
				assert.Containsf(t, svgFound, tt.svgMatch, bodyFound)
			}

			assert.NoErrorf(t, err, "chromedp run", bodyFound)
		})
	}

	t.Run("Jasmine", func(t *testing.T) {
		var jasmineOverallResult string

		err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("%s%s", e2eBaseUrl, "static/test.html")),
			chromedp.WaitVisible(`.jasmine-alert .jasmine-duration`),
			// // retrieve relevant values
			chromedp.Text(`.jasmine-alert .jasmine-overall-result`, &jasmineOverallResult),
		)

		assert.Contains(t, jasmineOverallResult, ", 0 failures")
		assert.NoErrorf(t, err, "chromedp run")
	})
}
