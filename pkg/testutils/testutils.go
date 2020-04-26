package testutils

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"github.com/peteraba/roadmapper/pkg/migrations"
	"github.com/peteraba/roadmapper/pkg/repository"
)

// if shouldUpdateGoldenFiles = true, golden files should be updated
var shouldUpdateGoldenFiles = flag.Bool("update", false, "update golden files")

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
