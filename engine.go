package gungnir

import (
	"net/http"
)

type Engine struct {
	tree *tree
}

func New(path string) *Engine {
	return &Engine{tree: newTree(path)}
}

func (e *Engine) Group(path string) *node {
	return e.tree.rootNode.Group(path)
}

func (e *Engine) Run(addr string, closeCh <-chan struct{}) error {
	srv := http.Server{Addr: addr, Handler: e}
	go func() {
		select {
		case <-closeCh: srv.Close()
		}
	}()

	return srv.ListenAndServe()
}

func (e *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h, pp := e.tree.matchHandler(req.URL.Path, req.Method)
	// NOT FOUND
	if h == nil {
		resp.WriteHeader(404)
		return
	}
	c := newCtx(resp, req, pp)
	h(c)
}
