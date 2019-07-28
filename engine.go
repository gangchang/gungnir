package gungnir

import (
	"net/http"
)

type Engine struct {
	root *route
}

func New(basePath string) *Engine {
	r := newRoute(basePath)
	return &Engine{
		root: r,
	}
}

func (e *Engine) Group(path string) *route {
	return e.root.Group(path)
}

func (e *Engine) Run(addr string, closeCh <-chan struct{}) error {
	srv := http.Server{Addr: addr, Handler: e}
	go func() {
		select {
		case <-closeCh:
			srv.Close()
		}
	}()

	return srv.ListenAndServe()
}

func (e *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// TODO add values
	c := newCtx(resp, req, nil)
	paths := getRequestPaths(req)
	r, _, pos := e.root.matchRoute(paths)
	if r == nil {
		resp.WriteHeader(404)
		return
	}
	handlerFn, _ := r.matchHandler(paths[pos:], req.Method)
	if handlerFn == nil {
		resp.WriteHeader(404)
		return
	}
	if r.doMiddleWareFns(c) {
		handlerFn(c)
	}
}
