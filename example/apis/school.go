package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gangchang/gungnir"
)

type School struct {
	ID     string
	Name   string
	Region string
}

func (s *School) CreateOne(c gungnir.Ctx) {
	if err := c.Bind(s); err != nil {
		c.SetCodeMsg(403, err.Error())
		return
	}
	data, _ := json.Marshal(s)
	if err := db.Put(s.generateKey(), data, nil); err != nil {
		c.SetCodeMsg(500, err.Error())
		return
	}

	c.SetCodeMsg(200, string(data))
	return
}

func (s *School) Path() string {
	return "school"
}

func (s *School) ReadOne(c gungnir.Ctx) {
	id, _ := c.GetPathParam("school_id")
	s.ID = id
	data, err := db.Get(s.generateKey(), nil)
	if err != nil {
		c.SetCodeMsg(403, err.Error())
		return
	}
	c.SetCodeMsg(200, data)
}

func (s *School) generateKey()  []byte {
	return []byte(fmt.Sprintf("%s/%s", s.Path(), s.ID))
}
