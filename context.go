package gungnir

import "net/http"

type Ctx struct {
	Resp http.Response
	Req *http.Request
}
