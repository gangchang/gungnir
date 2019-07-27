package main

import (
	"fmt"
	"github.com/n1ce37/gungnir"
)

func main() {
	e := gungnir.New("")
	userGroup := e.Group("user")
	{
		userGroup.GET("test", User)
		userGroup.GET("", User)
	}
	closeCh := make(chan struct{}, 1)
	e.Run(":3737", closeCh)
}

type UserT struct {
	Name string `json:"name"`
}

func User(ctx *gungnir.Ctx) {
	ut := &UserT{}
	ctx.Bind(ut)
	fmt.Println(ut)
	ctx.JSON(200, ut)
}
