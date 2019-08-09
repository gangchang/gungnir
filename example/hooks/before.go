package hooks

import (
	"github.com/gangchang/gungnir"
	"time"
)

func TimeBegin(c gungnir.Ctx) bool {
	c.Set("time", time.Now())
	return true
}
