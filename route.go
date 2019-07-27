package gungnir

import (
	"fmt"
	"net/http"
)

type handler func(*Ctx)

type Route struct {
	path urlPath
	subRoutes    []*Route

	handlers map[string][]handlerRoute

	middlerWareFns []middlerWareFn
}

type handlerRoute struct {
	urlPath
	handler handler
}

func NewRoute(path string) *Route{
	return &Route{path: newURLPath(path)}
}

func (r *Route) AddMillderWareFn(mwf middlerWareFn) {
	r.middlerWareFns = append(r.middlerWareFns, mwf)
}

func (r *Route) Group(path string) *Route {
	subRoute := &Route{
		path: newURLPath(path),
		middlerWareFns: r.middlerWareFns,
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
	//hr := handlerRoute{
	//	urlPath: newURLPath(path),
	//	handler: h,
	//}
	//r.postHandlers  = append(r.postHandlers, hr)
}

func (r *Route) findRoute(paths []string) (*Route, int) {
	nowRoute := r
	nowCnt := 0
	incCnt, matched := nowRoute.Match(paths)
	if matched==MatchedFull || matched==MatchedSub {
		nowCnt += incCnt
	} else {
		return nil, -1
	}
	for {
		matched, incCnt, matchRoute := nowRoute.MatchSub(paths[nowCnt+1:])
		nowCnt += incCnt
		switch matched {
		case MatchedFull:
			return matchRoute, nowCnt
		case MatchedSub:
			if len(r.subRoutes) == 0 {
				return matchRoute, nowCnt
			}
			nowRoute = matchRoute
			nowCnt += incCnt
			continue
		case MatchedNo:
			fmt.Println("not matched")
			return nil, -1
		default:
			return nil, -1
		}
	}
}

func (r *Route) findHandler(paths []string, method string) handler {
	methodHandlers, exists := r.handlers[method]
	if !exists {
		return nil
	}
	for _, v := range methodHandlers{
		if v.fullMatch(paths) {
			return v.handler
		}
	}
	return nil
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
		//r.do(r.getHandlers, paths)
	case http.MethodPost:
		//r.do(r.postHandlers, paths)
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
