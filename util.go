package gungnir

import (
	"strings"
)

func getPaths(path string) []string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	paths := strings.Split(path, "/")
	if len(paths) == 0 {
		panic("0 paths")
	}

	return paths
}

func getPathParams(paths []string) map[int]string {
	pathParams := make(map[int]string)
	for k, v := range paths {
		if strings.HasPrefix(v, ":") {
			pathParams[k] = v[1:]
		}
	}

	return pathParams
}

func getWildcard(path string) (string, string, bool) {
	if strings.HasPrefix(path, "{") && strings.Contains(path, "}") {
		pos := strings.Index(path, "}")
		key := path[:pos]
		kind := path[pos+1:]
		return key, kind, true
	}

	return path, "", false
}
