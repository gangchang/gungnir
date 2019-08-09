package hooks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gangchang/gungnir"
)

func TimeEnd(c gungnir.Ctx) bool {
	tb, _ := c.Get("time")
	tn := time.Now()
	fmt.Println(tn.Sub(tb.(time.Time)))
	return true
}

func Resp(c gungnir.Ctx) bool {
	code, msg := c.GetCodeMsg()
	msgByte, ok := msg.([]byte)
	if !ok {
		msgByte, _ = json.Marshal(msg)
	}
	c.WriteCode(code)
	//c.WriteHeader("Content-Type", "application/json")
	c.WriteHeader("Content-Type", "application/json")
	c.WriteHeader("test", "123")
	c.WriteBody(msgByte)
	return true
}
