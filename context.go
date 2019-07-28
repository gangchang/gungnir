package gungnir

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func newCtx(resp http.ResponseWriter, req *http.Request, pathValues map[string]string) *Ctx {
	return &Ctx{
		resp: resp,
		req:  req,
	}
}

type Ctx struct {
	resp http.ResponseWriter
	req  *http.Request
	pathValues map[string]string
}

// TODO now just json
func (c *Ctx) Bind(obj interface{}) error {
	data, _ := ioutil.ReadAll(c.req.Body)
	return json.Unmarshal(data, obj)
}

func (c *Ctx) Write (code int, obj interface{}) {
	data, _ := json.Marshal(obj)
	c.resp.Header().Add("Content-Type", "application/json")
	c.resp.WriteHeader(code)
	c.resp.Write(data)
}

func (c *Ctx) GetFromPath(key string) (string, bool) {
	value, exists := c.pathValues[key]
	return value, exists
}
