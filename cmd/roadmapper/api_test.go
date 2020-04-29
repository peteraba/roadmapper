// +build api

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	gofakeit "github.com/brianvoe/gofakeit/v5"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peteraba/roadmapper/pkg/roadmap"
)

const (
	baseUrl      = "http://localhost:1323/api"
	yamlFilePath = "../../api.yml"
	jsonFilePath = "../../api.json"
)

var (
	router     *openapi3filter.Router
	httpClient *http.Client

	minDate = time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC)
	maxDate = time.Date(2020, 4, 25, 0, 0, 0, 0, time.UTC)
)

// init will:
// - provide a seed for all used random generators
// - create a new router from the api.json file
//   but first it will check if api.yml is newer than api.json and re-converts it if needed
func init() {
	gofakeit.Seed(0)

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	if router != nil {
		return
	}

	yamlFile, err := os.Stat(yamlFilePath)
	if err != nil {
		panic("could not find '" + yamlFilePath + "'" + err.Error())
	}

	jsonFile, err := os.Stat(jsonFilePath)
	if err != nil || yamlFile.ModTime().After(jsonFile.ModTime()) {
		content, err := ioutil.ReadFile(yamlFilePath)
		if err != nil {
			panic("could not read '" + yamlFilePath + "': " + err.Error())
		}

		jsonContent, err := yaml.YAMLToJSON(content)
		if err != nil {
			panic("could not parse '" + yamlFilePath + "': " + err.Error())
		}

		err = ioutil.WriteFile(jsonFilePath, jsonContent, os.ModePerm)
		if err != nil {
			panic("could not write '" + yamlFilePath + "': " + err.Error())
		}

		time.Sleep(10 * time.Second)
	}

	router = openapi3filter.NewRouter().WithSwaggerFromFile(jsonFilePath)
}

// doHttpWithBody sends an HTTP request and returns an HTTP response and the body content
func doHttpWithBody(t *testing.T, req *http.Request, expectedStatusCode int) (*http.Response, []byte) {
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)

	require.NoError(t, err)
	resp.Body.Close()
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return resp, body
}

func newRoadmapPayload() roadmap.RoadmapExchange {
	p := gofakeit.Number(0, 20)
	m := gofakeit.Number(0, p)

	var (
		milestones []roadmap.Milestone
		projects   []roadmap.Project
		project    roadmap.Project
		ind        = 0
	)

	for i := 0; i < m; i++ {
		milestones = append(milestones, newMilestone())
	}

	for i := 0; i < p; i++ {
		project = newProject(m, ind)
		projects = append(projects, project)
		ind = nextIndentation(ind)
	}

	return roadmap.RoadmapExchange{
		Title:      newWords(),
		DateFormat: "2006-01-02",
		BaseURL:    gofakeit.URL(),
		Projects:   projects,
		Milestones: milestones,
	}
}

func newProject(milestoneCount, ind int) roadmap.Project {
	m := gofakeit.Number(0, milestoneCount)
	d := newDates()
	p := gofakeit.Number(0, 100)

	project := roadmap.Project{
		Indentation: uint8(ind),
		Title:       newWords(),
		Milestone:   uint8(m),
		Dates:       d,
		Percentage:  uint8(p),
	}

	return project
}

func newWords() string {
	var w []string

	for i := 0; i < gofakeit.Number(1, 5); i++ {
		w = append(w, gofakeit.HipsterWord())
	}

	return strings.Join(w, " ")
}

func nextIndentation(indentation int) int {
	return indentation - gofakeit.Number(-1, indentation)
}

func newDates() *roadmap.Dates {
	if gofakeit.Bool() {
		return nil
	}

	var (
		d0 = gofakeit.DateRange(minDate, maxDate)
		d1 = gofakeit.DateRange(minDate, maxDate)
	)

	if d0.Before(d1) {
		return &roadmap.Dates{
			StartAt: d0,
			EndAt:   d1,
		}
	}

	return &roadmap.Dates{
		StartAt: d1,
		EndAt:   d0,
	}
}

func getURLs() []string {
	var (
		urls []string
	)

	for i := 0; i < gofakeit.Number(0, 2); i++ {
		urls = append(urls, gofakeit.URL())
	}

	for i := 0; i < gofakeit.Number(0, 2); i++ {
		urls = append(urls, gofakeit.Word())
	}

	return urls
}

func newMilestone() roadmap.Milestone {
	return roadmap.Milestone{
		Title:      newWords(),
		DeadlineAt: newDateOptional(),
		URLs:       getURLs(),
	}
}

func newDateOptional() *time.Time {
	var (
		optionalDate *time.Time
	)

	if gofakeit.Bool() {
		date := gofakeit.DateRange(minDate, maxDate)
		optionalDate = &date
	}

	return optionalDate
}

func newCreateRoadmapRequest(t *testing.T, re roadmap.RoadmapExchange) *http.Request {
	marshaled, err := json.Marshal(re)
	require.NoError(t, err)

	url := fmt.Sprintf("%s/", baseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewReader(marshaled))
	require.NoError(t, err)

	req.Header.Add("Content-Type", `application/json`)

	return req
}

func TestApi_CreateRoadmap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Create request
		roadmapRequestData := newRoadmapPayload()
		req := newCreateRoadmapRequest(t, roadmapRequestData)

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
		resp, body := doHttpWithBody(t, req, http.StatusCreated)

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
