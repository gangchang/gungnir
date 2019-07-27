package gungnir

import (
	"fmt"
	"net/http"
)

type handler func(*Ctx)

type route struct {
	path urlPath
	subRoutes    []*route

	handlers map[string][]handlerRoute

	middlerWareFns []middlerWareFn
}

type handlerRoute struct {
	urlPath
	handler handler
}

func newRoute(path string) *route{
	return &route{
		path: newURLPath(path),
		handlers: make(map[string][]handlerRoute),
	}
}

func (r *route) AddMillderWareFn(mwf middlerWareFn) {
	r.middlerWareFns = append(r.middlerWareFns, mwf)
}

func (r *route) Group(path string) *route {
	subRoute := newRoute(path)
	subRoute.middlerWareFns = r.middlerWareFns

	r.subRoutes = append(r.subRoutes, subRoute)
	return subRoute
}

func (r *route) OPTIONS(path string, h handler) {
	r.addHandler(path, http.MethodOptions, h)
}

func (r *route) GET(path string, h handler) {
	r.addHandler(path, http.MethodGet, h)
}

func (r *route) HEAD(path string, h handler) {
	r.addHandler(path, http.MethodHead, h)
}

func (r *route) POST(path string, h handler) {
	r.addHandler(path, http.MethodPost, h)
}

func (r *route) PUT(path string, h handler) {
	r.addHandler(path, http.MethodPut, h)
}

func (r *route) PATCH(path string, h handler) {
	r.addHandler(path, http.MethodPatch, h)
}

func (r *route) DELETE(path string, h handler) {
	r.addHandler(path, http.MethodDelete, h)
}

func (r *route) TRACE(path string, h handler) {
	r.addHandler(path, http.MethodTrace, h)
}

func (r *route) CONNECT(path string, h handler) {
	r.addHandler(path, http.MethodConnect, h)
}

func (r *route) addHandler(path, method string, h handler) {
	hr := handlerRoute{
		urlPath: newURLPath(path),
		handler: h,
	}

	_, exists := r.handlers[method]
	if !exists{
		r.handlers[method] = make([]handlerRoute, 0)
	}

	r.handlers[method] = append(r.handlers[method], hr)
}




func (r *route) findRoute(paths []string) (*route, int) {
	nowRoute := r
	nowCnt := 0
	incCnt, matched := nowRoute.Match(paths)
	if matched==MatchedFull || matched==MatchedSub {
		nowCnt += incCnt
	} else {
		return nil, -1
	}
	for {
		matched, incCnt, matchRoute := nowRoute.MatchSub(paths[nowCnt:])
		nowCnt += incCnt
		switch matched {
		case MatchedFull:
			return matchRoute, nowCnt
		case MatchedSub:
			if len(matchRoute.subRoutes) == 0 {
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

func (r *route) findHandler(paths []string, method string) handler {
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

func (r *route) Match(paths []string) (int, Matched){
	cnt, matched := r.path.match(paths)
	switch matched {
	case MatchedFull:
		fallthrough
	case MatchedSub:
		return cnt, matched
	}

	return 0, MatchedNo
}

func (r *route) MatchSub(paths []string) (Matched, int, *route) {
	for _, v := range r.subRoutes {
		cnt, matched := v.path.match(paths)
		if matched == MatchedSub || matched == MatchedFull {
			return matched, cnt, v
		}
	}
	return MatchedNo, -1, nil
}

func (r *route) Do(paths []string, method string) {
	switch method {
	case http.MethodGet:
		//r.do(r.getHandlers, paths)
	case http.MethodPost:
		//r.do(r.postHandlers, paths)
	default:
		// not support method
	}
}

func (r *route) do(hrs []handlerRoute, paths []string) {
	for _, v := range hrs {
		if v.urlPath.fullMatch(paths) {
			v.handler(nil)
			return
		}
	}
	//404
}
