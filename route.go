package gungnir

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
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

// TODO check path
func (r *route) newSubRoute(path string) *route{
	return nil
}

// TODO
func (r *route) Install(obj interface{}) {
	typ := reflect.TypeOf(obj)
	//name := reflect.TypeOf(obj).Name()
names := strings.Split(typ.String(), ".")
name := strings.ToLower(names[len(names)-1])

	onePath := getOnePath(name)
	manyPath := getManyPath(name)

	if coh, ok := obj.(createOneHandler); ok {
		r.addHandler(onePath, http.MethodPost, coh.CreateOne)
	}
	if cmh, ok := obj.(createManyHandler); ok {
		r.addHandler(manyPath, http.MethodPost, cmh.CreateMany)
	}
	if goh, ok := obj.(readOneHandler); ok {
		r.addHandler(onePath, http.MethodGet, goh.ReadOne)
	}
	if gmh, ok := obj.(readManyHandler); ok {
		r.addHandler(manyPath, http.MethodGet, gmh.ReadMany)
	}
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
		matched, route, values, incCnt := nowRoute.match(paths[nowCnt:])
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
		subMatched, subValues, subCnt := v.path.matchRoute(paths)
		switch subMatched {
		case MatchedFull:
		case MatchedSub:
			if len(v.subRoutes) == 0 {
				matched = MatchedFull
			}
		}

		if matched == MatchedNo {
			continue
		}

		values = mergeMap(values, subValues)
		cnt += subCnt

			return matched, v, values, cnt
		}

	return MatchedNo, nil, nil, -1
}
