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
	path := purePath(req.URL.Path)
	return strings.Split(path, "/")
}
