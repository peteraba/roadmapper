// +build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/repository"
	"github.com/peteraba/roadmapper/pkg/roadmap"
	"github.com/peteraba/roadmapper/pkg/testutils"
)

func testApp(baseRepo repository.PgRepository, logger *zap.Logger, port uint, assetsDir string) func() {
	repo := roadmap.Repository{PgRepository: baseRepo}
	codeBuilder := code.Builder{}
	handler := newRoadmapHandler(logger, repo, codeBuilder, "", "", false)
	server := newServer(handler, assetsDir, "", "")
	teardown := server.StartWithTeardown(port)

	return teardown
}

func TestE2E_API(t *testing.T) {
	var (
		apiPort    uint = 9877
		apiBaseUrl      = "http://localhost:9877/"
		apiDbUser       = "rdmp"
		apiDbPass       = "rdmp"
		apiDbName       = "rdmp"
		minDate         = time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC)
		maxDate         = time.Date(2020, 4, 25, 0, 0, 0, 0, time.UTC)
	)

	// create a new logger
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// create a new database
	baseRepo, _, teardown := testutils.SetupRepository(t, "TestE2E_API", apiDbUser, apiDbPass, apiDbName, logger)
	defer teardown()

	// start up a new app
	appTeardown := testApp(baseRepo, logger, apiPort, "")
	defer appTeardown()

	httpClient := testutils.GetHttpClient()
	router := testutils.GetRouter(t)

	t.Run("success", func(t *testing.T) {
		// Create request
		roadmapRequestData := roadmap.NewRoadmapExchangeStub(0, 0, minDate, maxDate)
		req := newCreateRoadmapRequest(t, roadmapRequestData, apiBaseUrl)

		// Find route in the swagger file
		route, pathParams, err := router.FindRoute(req.Method, req.URL)
		require.NoError(t, err)

		// Validate request
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParams,
			Route:      route,
		}
		err = openapi3filter.ValidateRequest(nil, requestValidationInput)
		assert.NoError(t, err)

		// Get response
		resp, body := testutils.AssertHttp(t, httpClient, req, http.StatusCreated)

		// Validate response
		err = openapi3filter.ValidateResponse(nil, &openapi3filter.ResponseValidationInput{
			RequestValidationInput: requestValidationInput,
			Status:                 resp.StatusCode,
			Header:                 resp.Header,
			Body:                   resp.Body,
		})
		require.NoError(t, err)

		// Read and parse response
		response := roadmap.RoadmapExchange{}
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)
	})
}

func newCreateRoadmapRequest(t *testing.T, re roadmap.RoadmapExchange, baseUrl string) *http.Request {
	marshaled, err := json.Marshal(re)
	require.NoError(t, err)

	url := fmt.Sprintf("%s/api/", strings.TrimRight(baseUrl, "/"))
	req, err := http.NewRequest("POST", url, bytes.NewReader(marshaled))
	require.NoError(t, err)

	req.Header.Add("Content-Type", `application/json`)

	return req
}
