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
	if len(path) == 0 {
		return urlPath{
			Cnts: 0,
		}
	}
	paths := strings.Split(purePath(path), "/")
	return urlPath{
		Cnts: len(paths),
		paths: paths,
	}
}

func (up urlPath) match(paths []string) (int, Matched) {
	if up.Cnts == 0 {
		return 0, MatchedSub
	}
	if len(paths) < up.Cnts {
		return 0, MatchedNo
	}

	i := 0
	for i = range up.paths {
		if len(up.paths[i]) == 0 {
			continue
		}
		if strings.EqualFold(up.paths[i], paths[i]) {
			continue
		} else {
			return 0, MatchedNo
		}
	}
	i += 1

	if up.Cnts == len(paths) {
		return i, MatchedFull
	}

	return i, MatchedSub
}

func (up urlPath) fullMatch(paths []string) bool {
	if len(paths) == 0 {
		return true
	}
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

	return true
}
