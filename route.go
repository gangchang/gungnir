package gungnir

import (
	"fmt"
	"net/http"
	"strings"
)

type handler func(*Ctx)

type Route struct {
	path urlPath
	subRoutes    []*Route
	handlers map[string][]handlerRoute
}

type handlerRoute struct {
	urlPath
	handler handler
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
	hr := handlerRoute{
		urlPath: newURLPath(path),
		handler: h,
	}

	handlers, exists := r.handlers[http.MethodGet]
	if !exists {
		handlers := make([]handlerRoute, 1)
		r.handlers[http.MethodGet] = handlers
	}

	handlers = append(handlers, hr)
}

func (r *Route) POST(path string, h handler) {
	hr := handlerRoute{
		urlPath: newURLPath(path),
		handler: h,
	}
	r.postHandlers  = append(r.postHandlers, hr)
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
			if len(r.subRoutes) == 0 {
				return matchRoute
			}
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

func (r *Route) Do(paths []string, method string) {
	switch method {
	case http.MethodGet:
		r.do(r.getHandlers, paths)
	case http.MethodPost:
		r.do(r.postHandlers, paths)
	default:
		// not support method
	}
}

func (r *Route) do(hrs []handlerRoute, paths []string) {
	for _, v := range hrs {
		if v.urlPath.fullMatch(paths) {
			v.handler(nil)
			return
		}
	}
	//404
}
