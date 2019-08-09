package apis

import (
	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

func init() {
	var err error
	db, err = leveldb.OpenFile("db", nil)
	if err != nil {
		panic(err)
	}
}

func GetDB() *leveldb.DB {
	return db
}

func CloseDB() {
	db.Close()
}
