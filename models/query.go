package models

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/utils"
)

func open_array(vals []interface{}) []bson.ObjectId{
	result := make([]bson.ObjectId, 0)
	for _, val := range vals {
		result = append(result, val.(bson.ObjectId))
	}
	return result
}

func QueryDataIDWithMap(query map[string]interface{}) ([]bson.ObjectId, string) {
	db := S.DB("database")
	if query["key"] != nil && query["val"] != nil && query["op"] == nil && query["exps"] == nil {
		if query["key"] == "" || query["val"] == "" {
			var result bson.M
			err := db.C("table").Pipe([]bson.M{
				{"$group":bson.M{"_id":nil, "ids":bson.M{"$push":"$_id"}}},
			}).One(&result)
			//fmt.Println(len(result["ids"].([]interface{})))
			if err == nil {
				return open_array(result["ids"].([]interface{})), ""
			} else {
				return make([]bson.ObjectId, 0), "查询逻辑错误"
			}
		} else {
			var result bson.M
			err := db.C("table").Pipe([]bson.M{
				{"$match":bson.M{query["key"].(string)+".value":query["val"]}},
				{"$group":bson.M{"_id":nil, "ids":bson.M{"$push":"$_id"}}},
			}).One(&result)
			if err == nil {
				return open_array(result["ids"].([]interface{})), ""
			} else {
				return make([]bson.ObjectId, 0), "查询逻辑错误"
			}
		}
	} else if query["op"] != nil && query["exps"] != nil && query["key"] == nil && query["val"] == nil {
		if query["op"] == "and"  || query["op"] == "or" {
			result_set := utils.NewSimpleSet()
			err := ""
			arr := []bson.ObjectId{}
			for index, value := range query["exps"].([]interface{}) {
				arr, err = QueryDataIDWithMap(value.(map[string]interface{}))
				if err != "" {
					break
				}
				temp_set := utils.NewSimpleSetWithKeys(arr)
				if(query["op"] == "or") {
					result_set = utils.Union(result_set, temp_set)
				} else {
					if index == 0 {
						result_set = temp_set
					} else {
						result_set = utils.Intersect(result_set, temp_set)
					}
				}
			}
			if err == "" {
				return open_array(result_set.Elements()), ""
			} else {
				return make([]bson.ObjectId, 0), "查询逻辑错误"
			}
		} else {
			return make([]bson.ObjectId, 0), "查询逻辑错误"
		}
	} else {
		return make([]bson.ObjectId, 0), "查询逻辑错误"
	}
}

func QueryDataWithString(query string) map[string]interface{} {
	if(query == "") {
		m := GetAllData()
		m["success"] = false
		m["reason"] = "空查询"
		return m
	}
	var query_map map[string]interface{}
	err := json.Unmarshal([]byte(query), &query_map)
	if err != nil {
		m := GetAllData()
		m["success"] = false
		m["reason"] = "查询格式错误"
		return m
	}

	return QueryDataWithMap(query_map)
}

func QueryDataWithMap(query map[string]interface{}) map[string]interface{} {
	ids, err := QueryDataIDWithMap(query)
	if err == "" {
		return GetDataWithIDs(ids)
	} else {
		m := GetDataWithIDs(ids)
		if(m["success"] == nil && m["reason"] == nil) {
			m["success"] = false
			m["reason"] = err
		}
		return m
	}
}
