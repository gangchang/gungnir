package gungnir

import "strings"

type Matched int

const (
	MatchedFull Matched = iota
	MatchedSub
	MatchedNo
)

type urlPath struct{
	Cnts int
	paths []string
}

func newURLPath(path string) urlPath {
	paths := strings.Split(purePath(path), "/")
	return urlPath{
		Cnts: len(paths),
		paths: paths,
	}
}

func (up urlPath) match(paths []string) (int, Matched) {
	if len(paths) < up.Cnts {
		return -1, MatchedNo
	}

	i := 0
	for i = range up.paths {
		if len(up.paths[i]) == 0 {
			continue
		}
		if strings.EqualFold(up.paths[i], paths[i]) {
			continue
		} else {
			return -1, MatchedNo
		}
	}

	if up.Cnts == i+1 {
		// full match
		return i, MatchedFull
	}

	return i, MatchedSub
}

func (up urlPath) fullMatch(paths []string) bool {
	if len(paths) != up.Cnts {
		return false
	}
	i := 0
	for i = range up.paths {
		if len(up.paths[i]) == 0 {
			continue
		}
		if strings.EqualFold(up.paths[i], paths[i]) {
			continue
		} else {
			return false
		}
	}

	if i+1 != len(paths) {
		return false
	}

	return true
}
