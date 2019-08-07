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

func getWildcard(path string) (string, string, bool) {
	if strings.HasPrefix(path, "{") && strings.Contains(path, "}") {
		pos := strings.Index(path, "}")
		key := path[:pos]
		kind := path[pos+1:]
		return key, kind, true
	}

	return path, "", false
}

func mergeMap(map0 map[string]string, map1 map[string]string) map[string]string {
	if map0 == nil && map1 == nil {
		return nil
	} else if map0 == nil {
		return map1
	} else if map1 == nil {
		return map0
	}

	for k, v := range map1 {
		map0[k] = v
	}

	return map0
}


