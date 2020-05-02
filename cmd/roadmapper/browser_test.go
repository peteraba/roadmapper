// +build e2e

package main

import (
	"context"
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
	e2eTitle        = "Example Roadmap"
	e2eTxt          = `Monocle ipsum dolor sit amet
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
			svgMatch: "Monocle ipsum dolor sit",
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
}
