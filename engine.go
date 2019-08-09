package gungnir

import (
	"net/http"
)

type Engine struct {
	*tree
}

func New(path string) *Engine {
	return &Engine{tree: newTree(path)}
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
	n, h, pp := e.tree.matchHandler(req.URL.Path, req.Method)
	// NOT FOUND
	if h == nil {
		resp.WriteHeader(404)
		return
	}
	c := newCtx(resp, req, pp)
	for _, v := range n.beforeExecFns {
		if !v(c) {
			return
		}
	}
	h(c)
	for _, v := range n.afterExecFns {
		if !v(c) {
			return
		}
	}
}
