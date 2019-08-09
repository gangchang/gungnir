package main

import (
	"github.com/gangchang/gungnir"
	"github.com/gangchang/gungnir/example/apis"
	"github.com/gangchang/gungnir/example/hooks"
)

func main() {
	e := gungnir.New("")
	e.AddBeforeExecFns(hooks.TimeBegin)
	e.AddAfterExecFns(hooks.Resp, hooks.TimeEnd)
	e.Register(&apis.School{})
	e.Run(":3737", nil)
}
