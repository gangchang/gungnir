package gungnir

import (
	"log"
	"testing"
)

func TestRouteMatch(t *testing.T) {
	r := NewRoute("/aa")
	g := r.Group("bb")
	g.GET("bb", func() {})
	matched := r.Dispatch("/aa/bb", "")
	log.Println(matched)
	matched = r.Dispatch("/aa/cc", "")
	log.Println(matched)
}
