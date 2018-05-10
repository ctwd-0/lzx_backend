package models

import (
	"gopkg.in/mgo.v2"
)

//全局数据库链接
var S *mgo.Session

//建立数据库链接
func Start_server() {
	var err error
	S, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
}
