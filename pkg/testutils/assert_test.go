package testutils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"go.uber.org/zap"
)

func TestAssertQueries(t *testing.T) {
	type args struct {
		queries     []string
		assertLines []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"success - empty",
			args{
				[]string{},
				[]string{},
			},
		},
		{
			"success - non-empty",
			args{
				[]string{"SELECT 1"},
				[]string{"SELECT 1"},
			},
		},
		{
			"success - multiple",
			args{
				[]string{"SELECT 1", "SELECT 2"},
				[]string{"SELECT 1", "SELECT 2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, buf := SetupTestLogger()
			for _, q := range tt.args.queries {
				logger.Info("database query",
					zap.String("formattedQuery", q))
			}

			AssertQueries(t, buf, tt.args.assertLines)
		})
	}
}

func TestAssertQueriesRegexp(t *testing.T) {
	type args struct {
		queries     []string
		assertLines []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"success - empty",
			args{
				[]string{},
				[]string{},
			},
		},
		{
			"success - non-empty",
			args{
				[]string{"SELECT 1"},
				[]string{"SELECT 1"},
			},
		},
		{
			"success - multiple",
			args{
				[]string{"SELECT 1", "SELECT 2"},
				[]string{"SELECT 1", "SELECT 2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, buf := SetupTestLogger()
			for _, q := range tt.args.queries {
				logger.Info("database query",
					zap.String("formattedQuery", q))
			}

			AssertQueriesRegexp(t, buf, tt.args.assertLines)
		})
	}
}

type mockRoundTripper struct {
	*http.Response
}

func (mrt mockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return mrt.Response, nil
}

func TestAssertHttp(t *testing.T) {
	type args struct {
		req                *http.Request
		expectedStatusCode int
	}
	tests := []struct {
		name string
		args args
		body []byte
	}{
		{
			"abc",
			args{
				&http.Request{
					URL: &url.URL{},
				},
				200,
			},
			[]byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := ioutil.NopCloser(bytes.NewBuffer(tt.body))

			httpClient := &http.Client{}
			httpClient.Transport = mockRoundTripper{
				&http.Response{
					StatusCode: tt.args.expectedStatusCode,
					Body:       body,
				},
			}

			AssertHttp(t, httpClient, tt.args.req, tt.args.expectedStatusCode)
		})
	}
}

func TestAssertAPI(t *testing.T) {
	type args struct {
		yamlContent []byte
		jsonContent []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"success - empty",
			args{
				yamlContent: []byte(""),
				jsonContent: []byte("null"),
			},
		},
		{
			"success - non-empty",
			args{
				yamlContent: []byte(`
info:
  version: 0.0.1
  title: Roadmapper API
  description: API for Roadmapper
  termsOfService: https://docs.rdmp.app/terms/`),
				jsonContent: []byte(`{"info":{"version":"0.0.1","title":"Roadmapper API","description":"API for Roadmapper","termsOfService":"https://docs.rdmp.app/terms/"}}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertAPI(t, tt.args.yamlContent, tt.args.jsonContent)
		})
	}
}
