package gungnir

import (
	"fmt"
	"strings"
)

type handler func()

type Route struct {
	path urlPath
	subRoutes    []*Route
	getHandlers  []handler
	postHandlers []handler
}

func NewRoute(path string) *Route{
	return &Route{path: newURLPath(path)}
}

func (r *Route) Group(path string) *Route {
	subRoute := &Route{
		path: newURLPath(path),
	}
	r.subRoutes = append(r.subRoutes, subRoute)
	return subRoute
}

func (r *Route) GET(path string, h handler) {
	//if r.getHandlers == nil {
	//	r.getHandlers = make(map[urlPath]handler)
	//}
	//r.getHandlers[newURLPath(path)] = h
}

func (r *Route) POST(path string, h handler) {
	//if r.postHandlers == nil {
	//	r.postHandlers = make(map[urlPath]handler)
	//}
	//r.postHandlers[newURLPath(path)] = h
}

func (r *Route) Dispatch(path, method string) *Route {
	paths := strings.Split(purePath(path), "/")
	nowRoute := r
	nowCnt := 0
	increCnt, matched := nowRoute.Match(paths)
	if matched==MatchedFull || matched==MatchedSub {
		nowCnt += increCnt
	} else {
		return nil
	}
	for {
		matched, increCnt, matchRoute := nowRoute.MatchSub(paths[nowCnt+1:])
		switch matched {
		case MatchedFull:
			return matchRoute
		case MatchedSub:
			nowRoute = matchRoute
			nowCnt += increCnt
			continue
		case MatchedNo:
			fmt.Println("not matched")
			return nil
		default:
			return nil
		}
	}
}

func (r *Route) Match(paths []string) (int, Matched){
	cnt, matched := r.path.match(paths)
	switch matched {
	case MatchedFull:
		fallthrough
	case MatchedSub:
		return cnt, matched
	}

	return -1, MatchedNo
}

func (r *Route) MatchSub(paths []string) (Matched, int, *Route) {
	for _, v := range r.subRoutes {
		cnt, matched := v.path.match(paths)
		if matched == MatchedSub || matched == MatchedFull {
			return matched, cnt, v
		}
	}
	return MatchedNo, -1, nil
}
