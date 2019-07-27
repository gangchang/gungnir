package gungnir

import (
	"log"
	"net/http"
)

type Engine struct {
	root *Route
}

func (e *Engine) Run(addr string) {
	if err := http.ListenAndServe(addr, e); err != nil {
		log.Print(err)
	}
}

func (e *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := e.getContext(resp, req)
	paths := getRequestPaths(req)
	r, pos := e.root.findRoute(paths)
	if r == nil {
		resp.WriteHeader(404)
		return
	}
	handlerFn := r.findHandler(paths[pos:], req.Method)
	if handlerFn == nil {
		resp.WriteHeader(404)
		return
	}
	// TODO millderWareFns
	handlerFn(ctx)
}

func (e *Engine) getContext(resp http.ResponseWriter, req *http.Request)*Ctx {
	return newCtx(resp, req)
}
