package testutils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// init will:
// - provide a seed for all used random generators
// - create a new router from the api.json file
//   but first it will check if api.yml is newer than api.json and re-converts it if needed
func init() {
	gofakeit.Seed(0)
}

var (
	timeout time.Duration
)

func AssertQueries(t *testing.T, buf *zaptest.Buffer, lines []string) {
	logs := buf.Lines()
	if assert.Equal(t, len(lines), len(logs)) {
		for i, l := range lines {
			exp, err := json.Marshal(queryLog{l})
			require.NoError(t, err)

			assert.Equal(t, string(exp), logs[i])
		}
	}
}

func AssertQueriesRegexp(t *testing.T, buf *zaptest.Buffer, lines []string) {
	logs := buf.Lines()
	if assert.Equal(t, len(lines), len(logs)) {
		for i, l := range lines {
			assert.Regexp(t, l, logs[i])
		}
	}
}

// AssertHttp sends an HTTP request and returns an HTTP response and the body content
func AssertHttp(t *testing.T, httpClient *http.Client, req *http.Request, expectedStatusCode int) (*http.Response, []byte) {
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)

	require.NoError(t, err)
	resp.Body.Close()
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return resp, body
}

// AssertAPI compares the content of a .yml and .json formatted file, requires them to be the same
func AssertAPI(t *testing.T, yamlContent, jsonContent []byte) {
	yamlAsJson, err := yaml.YAMLToJSON(yamlContent)
	require.NoErrorf(t, err, "can't convert YAML to JSON")

	require.JSONEqf(t, string(yamlAsJson), string(jsonContent), "content of JSON and YAML APIs don't match")
}
