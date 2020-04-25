package testutils

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"
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
