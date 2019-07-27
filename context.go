package gungnir

import "net/http"

func newCtx(resp http.ResponseWriter, req *http.Request) *Ctx {
	return &Ctx{
		Resp: resp,
		Req: req,
	}
}

type Ctx struct {
	Resp http.ResponseWriter
	Req *http.Request
}

func (c *Ctx) Get() {

}
