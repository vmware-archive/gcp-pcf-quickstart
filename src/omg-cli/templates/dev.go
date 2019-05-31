// +build dev

package templates

import (
	"net/http"
	"path/filepath"
	"runtime"
)

func assetsDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "assets")
}

var Templates http.FileSystem = http.Dir(assetsDir())
