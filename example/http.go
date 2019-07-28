package main

import (
	"fmt"
	"github.com/n1ce37/gungnir"
)

func main() {
	e := gungnir.New("")
	userGroup := e.Group("user")
	{
		userGroup.POST("test", User)
		userGroup.GET("aa/bb", User)
		userGroup.GET("aa/cc/{id}", User)
	}
	teacherGroup := e.Group("school")
	teacherGroup.Install(&Teacher{})
	closeCh := make(chan struct{}, 1)
	e.Run(":3737", closeCh)
}

type UserT struct {
	Name string `json:"name"`
}

func User(c *gungnir.Ctx) {
	ut := &UserT{}
	c.Bind(ut)
	fmt.Println(ut)
	c.Write(200, ut)
}

type Teacher struct {

}

func (*Teacher) ReadOne(c *gungnir.Ctx) {
	c.Write(200, map[string]string{
		"msg": "readone",
	})
}

func (*Teacher) ReadMany(c *gungnir.Ctx) {
		c.Write(200, map[string]string{
		"msg": "readone",
	})
}
