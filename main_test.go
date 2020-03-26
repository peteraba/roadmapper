// +build integration

package main

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mafredri/cdp/protocol/dom"
	"github.com/pkg/errors"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

func TestIntegration_Browser(t *testing.T) {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		t.Fatalf("integration tests require a base url to run against")
	}

	formID := "roadmap-form"

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

	_, _, _ = formID, txt, txtBaseUrl

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, client, domContent, err := setup(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown(conn, domContent)

	t.Run("cdp test", func(t *testing.T) {
		err := openPage(ctx, client, domContent, baseUrl)
		if err != nil {
			t.Fatal(err)
		}

		assertUrl(t, ctx, client, baseUrl)
		assertHtml(t, ctx, client, formID)

		err = newScreenshot(ctx, client)
		if err != nil {
			t.Fatal(errors.Wrap(err, "failed to create a screenshot"))
		}
	})
}

func setup(ctx context.Context) (*rpcc.Conn, *cdp.Client, page.DOMContentEventFiredClient, error) {
	devtoolsUrl := os.Getenv("DEVTOOLS_URL")
	if devtoolsUrl == "" {
		devtoolsUrl = "http://127.0.0.1:9222"
	}

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New(devtoolsUrl)
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "can't reach DevTools at %s", devtoolsUrl)
		}
	}

	// Initiate a new RPC connection to the Chrome DevTools Protocol target.
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "!!! 2 !!!")
	}

	c := cdp.NewClient(conn)

	// Open a DOMContentEventFired client to buffer this event.
	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return conn, c, domContent, errors.Wrap(err, "!!! 3 !!!")
	}

	// Enable events on the Page domain, it's often preferable to create
	// event clients before enabling events so that we don't miss any.
	if err = c.Page.Enable(ctx); err != nil {
		return conn, c, domContent, errors.Wrap(err, "!!! 4 !!!")
	}

	return conn, c, domContent, nil
}

func teardown(conn *rpcc.Conn, domContent page.DOMContentEventFiredClient) {
	conn.Close()
	domContent.Close()
}

func openPage(ctx context.Context, c *cdp.Client, domContent page.DOMContentEventFiredClient, baseUrl string) error {
	// Create the Navigate arguments with the optional Referrer field set.
	navArgs := page.NewNavigateArgs(baseUrl)
	_, err := c.Page.Navigate(ctx, navArgs)
	if err != nil {
		return err
	}

	// Wait until we have a DOMContentEventFired event.
	if _, err = domContent.Recv(); err != nil {
		return err
	}

	return nil
}

func assertUrl(t *testing.T, ctx context.Context, c *cdp.Client, baseUrl string) {
	history, err := c.Page.GetNavigationHistory(ctx)
	if err != nil {
		t.Fatal("not able to retrieve the navigation history")
	}

	if history.Entries[history.CurrentIndex].URL != baseUrl {
		t.Fatalf("current url: %s, want: %s", history.Entries[history.CurrentIndex].URL, baseUrl)
	}
}

func assertHtml(t *testing.T, ctx context.Context, c *cdp.Client, verify string) {
	// Fetch the document root node. We can pass nil here
	// since this method only takes optional arguments.
	doc, err := c.DOM.GetDocument(ctx, nil)
	if err != nil {
		t.Fatal("not able to retrieve the document")
	}

	result, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &doc.Root.NodeID,
	})
	if err != nil {
		t.Fatal("retrieve the HTML root node")
	}

	if !strings.Contains(result.OuterHTML, verify) {
		t.Fatalf("outer HTML of root node does not contain: %s", verify)
	}
}

func newScreenshot(ctx context.Context, c *cdp.Client) error {
	// Capture a screenshot of the current page.
	screenshotName := "screenshot.jpg"
	screenshotArgs := page.NewCaptureScreenshotArgs().
		SetFormat("jpeg").
		SetQuality(80)
	screenshot, err := c.Page.CaptureScreenshot(ctx, screenshotArgs)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(screenshotName, screenshot.Data, 0644)
}
