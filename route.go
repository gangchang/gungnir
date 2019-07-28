package gungnir

import (
	"fmt"
	"net/http"
)

type handler func(*Ctx)

type route struct {
	path      urlPath
	subRoutes []*route

	handlers map[string][]handlerRoute

	middleWareFns []middleWareFn
}

type handlerRoute struct {
	urlPath
	handler handler
}

func newRoute(path string) *route {
	return &route{
		path:     newURLPath(path),
		handlers: make(map[string][]handlerRoute),
	}
}

func (r *route) AddMilldeWare(mwf middleWareFn) {
	r.middleWareFns = append(r.middleWareFns, mwf)
}

func (r *route) doMiddleWareFns(c *Ctx) bool {
	for _, fn := range r.middleWareFns {
		if !fn(c) {
			return false
		}
	}

	return true
}

func (r *route) Group(path string) *route {
	subRoute := newRoute(path)
	subRoute.middleWareFns = r.middleWareFns

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
	if !exists {
		r.handlers[method] = make([]handlerRoute, 0)
	}

	r.handlers[method] = append(r.handlers[method], hr)
}

func (r *route) matchRoute(paths []string) (*route, map[string]string, int) {
	nowRoute := r
	nowCnt := 0

	nowValues := make(map[string]string)
	for {
		matched, route, values, incCnt := nowRoute.match(paths)
		nowRoute = route
		nowCnt += incCnt
		nowValues = mergeMap(nowValues, values)

		switch matched {
		case MatchedFull:
			return nowRoute, nowValues, nowCnt
		case MatchedSub:
			continue
		case MatchedNo:
			fmt.Println("not matched")
			return nil, nil, -1
		}
	}
}

func (r *route) matchHandler(paths []string, method string) (handler, map[string]string) {
	methodHandlers, exists := r.handlers[method]
	if !exists {
		return nil, nil
	}

	for _, v := range methodHandlers {
		if values, ok := v.matchHandler(paths); ok {
			return v.handler, values
		}
	}

	return nil, nil
}

func (r *route) match(paths []string) (Matched, *route, map[string]string, int) {
	matched, values, cnt := r.path.matchRoute(paths)
	switch matched {
	case MatchedFull:
		return MatchedFull, r, values, cnt
	case MatchedNo:
		return MatchedNo, nil, nil, -1
	}

	paths = paths[cnt:]
	for _, v := range r.subRoutes {
		matched, subValues, subCnt := v.path.matchRoute(paths)
		if matched == MatchedFull || matched == MatchedSub {
			cnt += subCnt
			valuesMerged := mergeMap(values, subValues)
			return matched, v, valuesMerged, cnt
		}
	}

	return MatchedNo, nil, nil, -1
}
