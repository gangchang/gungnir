package gungnir

import (
	"fmt"
	"path/filepath"
)

type optionsHandler interface {
	OPTIONS(c *Ctx)
}

type getHandler interface {
	Get(c *Ctx)
}

type GetCollectionHandler interface {
	GetCollection(c *Ctx)
}

type headHandler interface {
	HEAD(c *Ctx)
}

type postHandler interface {
	POST(c *Ctx)
}

type putHandler interface {
	PUT(c *Ctx)
}

type patchHandler interface {
	PATCH(c *Ctx)
}

type deleteHandler interface {
	DELETE(c *Ctx)
}

type createOneHandler interface {
	CreateOne(c *Ctx)
}

type createManyHandler interface {
	CreateMany(c *Ctx)
}

type readOneHandler interface {
	ReadOne(c *Ctx)
}

type readManyHandler interface {
	ReadMany(c *Ctx)
}

func getOnePath(name string) string {
	return filepath.Join(getManyPath(name), getID(name))
}

func getManyPath(name string) string {
	return fmt.Sprintf("%ss", name)
}

func getID(name string) string {
	return fmt.Sprintf("%s_id", name)
}
