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
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	SetCodeMsg(code int, msg interface{})
	GetCodeMsg() (int, interface{})

	WriteHeader(key, value string)
	WriteBody(body []byte)
	WriteCode(code int)
}

func newCtx(resp http.ResponseWriter, req *http.Request, pp map[string]string) Ctx {
	return &ctx{
		resp: resp,
		req:  req,
		pp:   pp,
	}
}

type ctx struct {
	resp http.ResponseWriter
	req  *http.Request
	pp   map[string]string
	data map[string]interface{}

	code int
	msg  interface{}
}

// TODO now just json
func (c *ctx) Bind(obj interface{}) error {
	data, _ := ioutil.ReadAll(c.req.Body)
	return json.Unmarshal(data, obj)
}

func (c *ctx) Write(code int, obj interface{}) {
	data, ok := obj.([]byte)
	if !ok {
		data, _ = json.Marshal(obj)
	}
	c.resp.WriteHeader(code)
	c.resp.Header().Add("Content-Type", "application/json")
	c.resp.Write(data)
}

func (c *ctx) WriteHeader(key, value string) {
	//c.resp.Header().Add(key, value)
	c.resp.Header().Set(key, value)
}

func (c *ctx) WriteBody(body []byte)  {
	c.resp.Write(body)
}

func (c *ctx) WriteCode(code int) {
	c.resp.WriteHeader(code)
}

func (c *ctx) GetPathParam(key string) (string, bool) {
	value, exists := c.pp[key]
	return value, exists
}

func (c *ctx) Set(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
}

func (c *ctx) Get(key string) (interface{}, bool) {
	value, exists := c.data[key]
	return value, exists
}

func (c *ctx) SetCodeMsg(code int, msg interface{}) {
	c.code = code
	c.msg = msg
}

func (c *ctx) GetCodeMsg() (int, interface{}) {
	return c.code, c.msg
}
