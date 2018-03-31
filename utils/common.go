package utils

import (
	"gopkg.in/mgo.v2/bson"
)

func StringToData(new_folder [][]string) []bson.M {
	data := []bson.M{}
	for idx, val := range new_folder[0] {
		data = append(data, bson.M{"val":val,"key":bson.ObjectIdHex(new_folder[1][idx])})
	}
	return data
}
