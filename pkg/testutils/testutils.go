package testutils

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/ghodss/yaml"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"github.com/peteraba/roadmapper/pkg/migrations"
	"github.com/peteraba/roadmapper/pkg/repository"
)

var (
	// if shouldUpdateGoldenFiles = true, golden files should be updated
	shouldUpdateGoldenFiles = flag.Bool("update", false, "update golden files")

	testRouter     *openapi3filter.Router
	testHttpClient *http.Client

	yamlFilePath = "../../api.yml"
	jsonFilePath = "../../api.json"
)

// ShouldUpdateGoldenFiles return true if golden files should be updated
func ShouldUpdateGoldenFiles() bool {
	return shouldUpdateGoldenFiles != nil && *shouldUpdateGoldenFiles
}

func resourceFilePath(pathParts ...string) string {
	_, filename, _, _ := runtime.Caller(0)
	parentDir, _ := path.Split(path.Dir(filename))
	return path.Join(parentDir, "..", "res", path.Join(pathParts...))
}

func LoadFile(t *testing.T, pathParts ...string) []byte {
	filePath := resourceFilePath(pathParts...)

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	return dat
}

func SaveFile(t *testing.T, content []byte, pathParts ...string) {
	filePath := resourceFilePath(pathParts...)

	err := ioutil.WriteFile(filePath, content, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

// SetupRepository returns a new repository and a tear down function
func SetupRepository(t *testing.T, appName, dbUser, dbPass, dbName string, logger *zap.Logger, fixture ...interface{}) (repository.PgRepository, func(fixture ...interface{}), func()) {
	pool, resource, port := setupDb(t, dbUser, dbPass, dbName)

	repo := repository.NewPgRepository(appName, "localhost", port, dbName, dbUser, dbPass, logger)

	conn := repo.ConnectNoHook()
	m := migrations.New(dbUser, dbPass, "localhost", port, dbName)

	down := func() {
		_, err := m.Down(0)
		require.NoErrorf(t, err, "failed to run down migrations")
	}

	up := func(fixture ...interface{}) {
		_, err := m.Up(0)
		require.NoErrorf(t, err, "failed to run migrations")

		for i := range fixture {
			res, err := conn.Model(fixture[i]).Insert()
			require.NoErrorf(t, err, "failed to insert fixture", i)
			require.Equalf(t, 1, res.RowsAffected(), "no rows affected")
		}
	}
	up(fixture...)

	reset := func(fixture ...interface{}) {
		down()
		up(fixture...)
	}

	teardown := func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("could not tear down the database: %v", err)
		}
	}

	return repo, reset, teardown
}

func setupDb(t *testing.T, dbUser, dbPass, dbName string) (*dockertest.Pool, *dockertest.Resource, string) {
	var db *sql.DB
	var err error

	dbPool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	dbResource, err := dbPool.Run("postgres", "alpine", []string{"POSTGRES_USER=" + dbUser, "POSTGRES_PASSWORD=" + dbPass, "POSTGRES_DB=" + dbName})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	if err = dbPool.Retry(func() error {
		var err error
		dataSource := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, dbResource.GetPort("5432/tcp"), dbName)
		db, err = sql.Open("postgres", dataSource)
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	dbPort := dbResource.GetPort("5432/tcp")

	return dbPool, dbResource, dbPort
}

func SetupTestLogger() (*zap.Logger, *zaptest.Buffer) {
	buf := &zaptest.Buffer{}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{}),
		buf,
		zap.DebugLevel,
	))

	return logger, buf
}

type queryLog struct {
	FormattedQuery string `json:"formattedQuery"`
}

func GetRouter(t *testing.T) *openapi3filter.Router {
	if testRouter != nil {
		return testRouter
	}

	yamlFile, err := os.Stat(yamlFilePath)
	require.NoError(t, err, "can't stat .yml file: ", yamlFile)

	jsonFile, err := os.Stat(jsonFilePath)
	if err != nil || yamlFile.ModTime().After(jsonFile.ModTime()) {
		updateJson(t)
	} else {
		verifyJson(t)
	}

	testRouter = openapi3filter.NewRouter().WithSwaggerFromFile(jsonFilePath)

	return testRouter
}

func updateJson(t *testing.T) {
	content, err := ioutil.ReadFile(yamlFilePath)
	require.NoError(t, err, "could not read .yml file: ", yamlFilePath)

	jsonContent, err := yaml.YAMLToJSON(content)
	require.NoError(t, err, "could not parse .yml file: ", yamlFilePath)

	err = ioutil.WriteFile(jsonFilePath, jsonContent, os.ModePerm)
	require.NoError(t, err, "could not write .json file: ", jsonFilePath)
}

func verifyJson(t *testing.T) {
	yamlContent, err := ioutil.ReadFile(yamlFilePath)
	require.NoError(t, err, "could not read .yml file: ", yamlFilePath)

	jsonContent, err := ioutil.ReadFile(jsonFilePath)
	require.NoError(t, err, "could not read .json file: ", jsonFilePath)

	AssertAPI(t, yamlContent, jsonContent)
}

func GetHttpClient() *http.Client {
	if testHttpClient == nil {
		testHttpClient = &http.Client{}
	}

	return testHttpClient
}

func GetTimeout(t *testing.T) time.Duration {
	if timeout != 0 {
		return timeout
	}

	timeout = 30 * time.Second
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
