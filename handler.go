package gungnir

import (
	"fmt"
)

type typePath interface {
	Path() string
}

type createOneHandler interface {
	CreateOne(c Ctx)
}

type createManyHandler interface {
	CreateMany(c Ctx)
}

type readOneHandler interface {
	ReadOne(c Ctx)
}

type readManyHandler interface {
	ReadMany(c Ctx)
}

type deleteOneHandler interface {
	DeleteOne(c Ctx)
}

type deleteManyHandler interface{
	DeleteMany(c Ctx)
}

type updateOneHandler interface {
	UpdateOne(c Ctx)
}

type updateManyHandler interface {
	UpdateMany(c Ctx)
}

type optionsHandler interface {
	OPTIONS(c Ctx)
}

type getHandler interface {
	Get(c Ctx)
}

type GetCollectionHandler interface {
	GetCollection(c Ctx)
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

func getIDPath(name string) string {
	return fmt.Sprintf(":%s_id", name)
}

func getCollectionPath(name string) string {
	return fmt.Sprintf("%ss", name)
}
