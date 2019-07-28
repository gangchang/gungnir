package gungnir

import (
	"strings"
)

type Matched int

const (
	MatchedFull Matched = iota
	MatchedSub
	MatchedNo
)

type urlPath struct {
	cnt         int
	sections    []section
	wildcardPos map[int]struct{}
}

type section struct {
	path string
}

func newSection(path string) (section, bool) {
	sp, _, ok := getWildcard(path)

	return section{path:sp}, ok
}

func newURLPath(path string) urlPath {
	if len(path) == 0 {
		return urlPath{
			cnt: 0,
		}
	}

	paths := strings.Split(purePath(path), "/")
	up := urlPath{
		cnt:         len(paths),
		wildcardPos: make(map[int]struct{}),
	}
	for i, v := range paths {
		s, isWildcard := newSection(v)
		up.sections = append(up.sections, s)
		if isWildcard {
			up.wildcardPos[i] = struct{}{}
		}
	}

	return up
}

func (up urlPath) matchRoute(paths []string) (Matched, map[string]string, int) {
	if up.cnt == 0 {
		return MatchedSub, nil, 0
	}
	if len(paths) < up.cnt {
		return MatchedNo, nil, -1
	}

	length, ok := up.match(paths)
	if !ok {
		return MatchedNo, nil, -1
	}

	length += 1
	matched := MatchedSub
	if up.cnt == len(paths) {
		matched = MatchedFull
	}
	values := up.getPathValues(paths)

	return matched, values, length
}

func (up urlPath) matchHandler(paths []string) (map[string]string, bool) {
	if len(paths) == 0 && len(up.sections) == 0 {
		return nil, true
	}
	if len(paths) != up.cnt {
		return nil, false
	}

	_, ok := up.match(paths)
	if !ok {
		return nil, false
	}

	values := up.getPathValues(paths)

	return values, true
}

func (up urlPath) match(paths []string) (int, bool) {
	i := 0
	for i=0; i<len(up.sections);i++ {
		if _, exists := up.wildcardPos[i]; exists {
			continue
		}
		if strings.EqualFold(up.sections[i].path, paths[i]) {
			continue
		} else {
			return -1, false
		}
	}

	return i, true
}

func (up urlPath) getPathValues(paths []string) map[string]string {
	wps := len(up.wildcardPos)
	if wps == 0 {
		return nil
	}

	values := make(map[string]string, wps)
		for k := range up.wildcardPos {
		values[up.sections[k].path] = paths[k]
	}

		return values
}

