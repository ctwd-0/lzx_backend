package models

import (
	"fmt"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/utils"
	"io/ioutil"
	"strings"
)

const preprocess_dest_dir = "E:/20170109/building_viewer/dist/files/"
var S *mgo.Session

func Start_server() {
	var err error
	S, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
}

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
			err := db.C("data").Pipe([]bson.M{
				{"$group":bson.M{"_id":nil, "ids":bson.M{"$push":"$_id"}}},
			}).One(&result)
			fmt.Println(len(result["ids"].([]interface{})))
			if err == nil {
				return open_array(result["ids"].([]interface{})), ""
			} else {
				return make([]bson.ObjectId, 0), "查询逻辑错误"
			}
		} else {
			var result bson.M
			err := db.C("data").Pipe([]bson.M{
				{"$match":bson.M{query["key"].(string)+".cur":query["val"]}},
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
			for index, value := range query["exps"].([]interface{}) {
				arr, err := QueryDataIDWithMap(value.(map[string]interface{}))
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

func GetDataWithIDs(ids []bson.ObjectId) map[string]interface{} {
	db := S.DB("database")
	header, selector, err := GetDataHeaderAndSelector()
	var data []interface{}
	if err == nil {
		err = db.C("data").Find(bson.M{"_id":bson.M{"$in":ids}}).Select(selector).All(&data)
	}
	m := make(map[string]interface{})
	if err != nil {
		m["success"] = false
		m["reason"] = "db error"
	} else {
		content := convertData(data, header)
		m["header"] = header
		m["content"] = content
	}
	return m
}

func GetDataHeaderAndSelector() (interface{}, bson.M, error) {
	db := S.DB("database")
	var header interface{}
	err := db.C("column_name").Find(bson.M{}).Select(bson.M{"_id":0,"cur":1}).One(&header)
	if err == nil {
		selector := make(bson.M)
		selector["_id"] = 0
		for _, value := range header.(bson.M)["cur"].([]interface{}) {
		 	selector[value.(string) + ".old"] = 0
		}
		return header.(bson.M)["cur"], selector, nil
	} else {
		return nil, nil, err
	}
}

func GetAllData() map[string]interface{} {
	db := S.DB("database")
	header, selector, err := GetDataHeaderAndSelector()
	var data []interface{}
	if err == nil {
		err = db.C("data").Find(bson.M{}).Select(selector).All(&data)
	}

	m := make(map[string]interface{})

	if err != nil {
		m["success"] = false
		m["reason"] = "db error"
	} else {
		content := convertData(data, header)
		m["header"] = header
		m["content"] = content
	}
	return m
}

func convertData(data []interface{}, header interface{}) [][]string {
	content := make([][]string, 0)
	for _, value := range data {
		var line []string
		for _, hd := range header.([]interface{}) {
			line = append(line, value.(bson.M)[hd.(string)].(bson.M)["cur"].(string))
		}
		content = append(content, line)
	}
	return content
}

func InitDbFile(path string) {
	groups, _ := ioutil.ReadDir(path)
	db := S.DB("database")
	c := db.C("file")

	for _, group := range groups {
		categories, _ := ioutil.ReadDir(path + "/" + group.Name())
		for _, cat := range categories {
			files, _ := ioutil.ReadDir(path + "/" + group.Name() + "/" + cat.Name())
			for _, file := range files {
				model_id := group.Name()
				category := cat.Name()
				file_name := file.Name()
				path_name := path + "/" + model_id + "/" + category + "/"+ file_name
				parts := strings.Split(file_name, ".")
				if parts[len(parts) - 2] == "thumbnail" {
					continue
				}
				fmt.Print(path_name)
				m, _ := utils.Thumbnail(path_name)
				if m["success"].(bool) {
					ext := m["ext"].(string)
					ori_md5 := utils.FileMD5(path_name)
					thumbnail_md5 := utils.FileMD5(m["thumbnail_path"].(string))
					utils.CopyFile(m["ori_path"].(string),
						preprocess_dest_dir + ori_md5 + "." + ext)
					utils.CopyFile(m["thumbnail_path"].(string),
						preprocess_dest_dir + thumbnail_md5 + "." + ext)
					if m["success"].(bool) {
						c.Insert(bson.M{
							"model_id": model_id,
							"category": category,
							"original_md5": ori_md5,
							"thumbnail_md5": thumbnail_md5,
							"original_saved_as": ori_md5 + "." + ext,
							"thumbnail_saved_as": thumbnail_md5 + "." + ext,
							"original_path": "/dist/files/" + ori_md5 + "." + ext,
							"thumbnail_path": "/dist/files/" + thumbnail_md5 + "." + ext,
							"original_name": file_name,
							"type": "image",
							"deleted": false,
						})
					}
				}
				fmt.Println(" finished")
			}
		}
	}
}
