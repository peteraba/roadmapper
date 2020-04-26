package testutils

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/ory/dockertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
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

func SetupDb(t *testing.T, dbUser, dbPass, dbName string) (*dockertest.Pool, *dockertest.Resource, string) {
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
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, dbResource.GetPort("5432/tcp"), dbName))
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

func TeardownDb(t *testing.T, dbPool *dockertest.Pool, dbResource *dockertest.Resource) {
	if err := dbPool.Purge(dbResource); err != nil {
		t.Fatalf("Could not tear down the database: %s", err)
	}
}

func SetupLogger(t *testing.T) (*zap.Logger, *zaptest.Buffer) {
	buf := &zaptest.Buffer{}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{}),
		buf,
		zap.DebugLevel,
	))

	return logger, buf
}
