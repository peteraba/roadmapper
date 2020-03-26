// +build integration

package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Base(t *testing.T) {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		t.Fatalf("integration tests require a base url to run against")
	}

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

	txtBaseUrl := "https://github.com/peteraba/roadmapper"

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	t.Run("chromedp test", func(t *testing.T) {
		var txtFound, txtBaseUrlFound, svgFound string

		err := chromedp.Run(ctx,
			chromedp.Navigate(baseUrl),
			// wait for form element to become visible (ie, page is loaded)
			chromedp.WaitVisible(`#roadmap-form`),
			// set the value of the textarea
			chromedp.SetValue(`#txt`, txt),
			// set the value of the base url
			chromedp.SetValue(`#base-url`, txtBaseUrl),
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
			t.Fatalf("chromedp test: error = %v, wantErr %v", err, false)
		}

		assert.Equal(t, txt, txtFound)
		assert.Equal(t, txtBaseUrl, txtBaseUrlFound)
		assert.Contains(t, svgFound, "Initial development")
	})
}
