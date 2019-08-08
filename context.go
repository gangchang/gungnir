package gungnir

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Ctx interface {
	Bind(obj interface{}) error
	Write(code int, obj interface{})
	GetPathParam(key string) (string, bool)
}

func newCtx(resp http.ResponseWriter, req *http.Request, pp map[string]string) Ctx {
	return &ctx{
		resp: resp,
		req:  req,
		pp: pp,
	}
}

type ctx struct {
	resp       http.ResponseWriter
	req        *http.Request
	pp map[string]string
}

// TODO now just json
func (c *ctx) Bind(obj interface{}) error {
	data, _ := ioutil.ReadAll(c.req.Body)
	return json.Unmarshal(data, obj)
}

func (c *ctx) Write(code int, obj interface{}) {
	data, _ := json.Marshal(obj)
	c.resp.Header().Add("Content-Type", "application/json")
	c.resp.WriteHeader(code)
	c.resp.Write(data)
}

func (c *ctx) GetPathParam(key string) (string, bool) {
	value, exists := c.pp[key]
	return value, exists
}
