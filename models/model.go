package models

import (
	"gopkg.in/mgo.v2"
)

var S *mgo.Session

func Start_server() {
	var err error
	S, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
}
