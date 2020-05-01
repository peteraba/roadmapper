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

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"github.com/peteraba/roadmapper/pkg/bindata"
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
func SetupRepository(t *testing.T, appName, dbUser, dbPass, dbName string, logger *zap.Logger, fixture ...interface{}) (repository.PgRepository, func()) {
	pool, resource, port := setupDb(t, dbUser, dbPass, dbName)

	repo := repository.NewPgRepository(appName, "localhost", port, dbName, dbUser, dbPass, logger)

	m := migrations.New(dbUser, dbPass, "localhost", port, dbName)
	_, err := m.Up(0)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	conn := repo.ConnectNoHook()
	for i := range fixture {
		res, err := conn.Model(fixture[i]).Insert()
		require.NoError(t, err)
		require.Equal(t, 1, res.RowsAffected())
	}

	return repo, func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("could not tear down the database: %v", err)
		}
	}
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

func SetupLogger() (*zap.Logger, *zaptest.Buffer) {
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

	yamlContent, err := bindata.Asset(yamlFilePath)
	require.NoErrorf(t, err, "unable to load YAML API content")

	jsonContent, err := bindata.Asset(jsonFilePath)
	require.NoErrorf(t, err, "unable to load JSON API content")

	AssertAPI(t, yamlContent, jsonContent)

	tr := openapi3filter.NewRouter()
	sl := openapi3.NewSwaggerLoader()
	s, err := sl.LoadSwaggerFromData(jsonContent)
	require.NoErrorf(t, err, "invalid JSON API content")

	err = tr.AddSwagger(s)
	require.NoErrorf(t, err, "can't add swagger to router")

	testRouter = tr

	return testRouter
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

	timeout = 15 * time.Second
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
