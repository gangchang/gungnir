package gungnir

import "strings"

type urlPath struct{
	Cnts int
	paths []string
}

func newURLPath(path string) urlPath {
	paths := strings.Split(path, "/")
	return urlPath{
		Cnts: len(paths),
		paths: paths,
	}
}

func (up urlPath) match(paths []string) ([]string, bool) {
	if len(paths) < up.Cnts {
		return nil, false
	}

	for i := range up.paths {
		// *
		if len(up.paths[i]) == 0 {
			continue
		}
		if strings.EqualFold(up.paths[i], paths[i]) {
			continue
		} else {
			return nil, false
		}
	}

	return paths[up.Cnts:], true
}
