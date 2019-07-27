package gungnir

import (
	"net/http"
	"strings"
)

func purePath(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return path
}

func getRequestPaths(req *http.Request) []string {
	path := req.URL.RawPath
	path = purePath(path)
	return strings.Split(path, "/")
}
