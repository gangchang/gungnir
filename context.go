package gungnir

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

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

// TODO now just json
func (c *Ctx) Bind(obj interface{}) error {
	data, _ := ioutil.ReadAll(c.Req.Body)
	return json.Unmarshal(data, obj)
}

func (c *Ctx) JSON(code int, obj interface{}) {
	data, _ := json.Marshal(obj)
	c.Resp.Header().Add("Content-Type", "application/json")
	c.Resp.WriteHeader(code)
	c.Resp.Write(data)
}
