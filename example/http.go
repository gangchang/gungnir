package main

import (
	"fmt"
	"github.com/n1ce37/gungnir"
)

func main() {
	e := gungnir.New("")
	users := e.Group("users/123")
	{
		users.GET(User)
		user := users.Group(":id/234")
		user.GET(UserFind)
	}
	closeCh := make(chan struct{}, 1)
	e.Run(":3737", closeCh)
}

type UserT struct {
	Name string `json:"name"`
	ID string `json:"id"`
}

func User(c gungnir.Ctx) {
	ut := &UserT{}
	c.Bind(ut)
	fmt.Println(ut)
	c.Write(200, ut)
}

func UserFind(c gungnir.Ctx) {
	ut := &UserT{}
	c.Bind(ut)
	ut.ID, _ = c.GetPathParam("id")
	fmt.Println(ut)
	c.Write(200, ut)
}

//type Teacher struct {
//}
//
//func (*Teacher) ReadOne(c *gungnir.Ctx) {
//	c.Write(200, map[string]string{
//		"msg": "readone",
//	})
//}
//
//func (*Teacher) ReadMany(c *gungnir.Ctx) {
//	c.Write(200, map[string]string{
//		"msg": "readone",
//	})
//}
//