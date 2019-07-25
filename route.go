package gungnir

import "strings"

type handler func()

type Route struct {
	path urlPath
	subRoutes    map[urlPath]*Route
	getHandlers  map[urlPath]handler
	postHandlers map[urlPath]handler
}

func (r *Route) Group(path string) *Route {
	subRoute := &Route{}
	r.subRoutes[newURLPath(path)] = subRoute
	return subRoute
}

func (r *Route) GET(path string, h handler) {
	if r.getHandlers == nil {
		r.getHandlers = make(map[urlPath]handler)
	}
	r.getHandlers[newURLPath(path)] = h
}

func (r *Route) POST(path string, h handler) {
	if r.postHandlers == nil {
		r.postHandlers = make(map[urlPath]handler)
	}
	r.postHandlers[newURLPath(path)] = h
}

func (r *Route) Dispatch(path, method string) {
	paths := strings.Split(path, "/")
	r.Match(paths)
}

func (r *Route) Match(paths []string) ([]string, bool, bool){
	subPaths, matched := r.path.match(paths)
	if !matched {
		return nil,  false, false
	}

	if len(subPaths) == 0 {
		return nil,  true, true
	}

	return subPaths, true, false
}

func (r *Route) MatchSub(paths []string) ([]string, *Route, bool) {
	return nil, nil, false
}

