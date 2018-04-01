package models

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"time"
	"gopkg.in/mgo.v2/bson"
)

func InitDbTable(path string) {
	db := S.DB("database")
	reason := ""
	file, err := os.Open(path)
	if err != nil {
		reason = "打开文件失败"
	} else {
		defer file.Close()
	}

	var body []byte
	if reason == "" {
		body, err = ioutil.ReadAll(file)
		if err != nil {
			reason = "读取文件失败"
		}
	}

	if reason == "" {
		str := string(body)
		parts := strings.Split(str, "\r\n")
		headers := strings.Split(parts[0], ",")
		contents := parts[1:]
		hdkv := []bson.M{}
		for _, val := range headers{
			hdkv = append(hdkv, bson.M{"val":val,"key":bson.NewObjectId()})
		}
		db.C("column").Insert(bson.M{
			"value": hdkv,
			"author": "system_init",
			"modified": time.Now(),
			"old": []bson.M{},
			"deleted":false,
		})

		for _, content := range contents {
			vals := strings.Split(content, ",")
			m := bson.M{}
			m["deleted"] = false
			if len(vals) == len(headers) {
				for index, value := range vals {
					m[hdkv[index]["key"].(bson.ObjectId).Hex()] = bson.M{
						"value": value,
						"modified": time.Now(),
						"author": "system_init",
						"old": []bson.M{},
					}
				}
				db.C("table").Insert(m)
			}
		}
	}

	fmt.Println("reason:", reason)
}

func GetDataWithIDs(ids []bson.ObjectId) map[string]interface{} {
	db := S.DB("database")
	header, selector, err := GetDataHeaderAndSelector()
	var data []bson.M
	if err == nil {
		err = db.C("table").Find(bson.M{"_id":bson.M{"$in":ids}}).Select(selector).All(&data)
	}
	m := make(map[string]interface{})
	if err != nil {
		m["success"] = false
		m["reason"] = "数据库错误"
	} else {
		content, ids := convertData(data, header)
		m["header"] = header[0]
		m["content"] = content
		m["ids"] = ids
	}
	return m
}

func GetDataHeader() ([]string, string) {
	db := S.DB("database")
	var header bson.M
	err := db.C("column").Find(bson.M{}).Select(bson.M{"_id":0,"value":1}).One(&header)
	if err == nil {
		hds := []string{}
		for _, value := range header["value"].([]interface{}) {
			val := value.(bson.M)["val"].(string)
			hds = append(hds, val)
		}
		return hds, ""
	} else {
		return []string{}, "数据库错误"
	}
}

func GetDataHeaderAndSelector() ([][]string, bson.M, error) {
	db := S.DB("database")
	var header bson.M
	err := db.C("column").Find(bson.M{}).Select(bson.M{"_id":0,"value":1}).One(&header)
	if err == nil {
		selector := bson.M{}
		//selector["_id"] = 0
		hds := [][]string{[]string{}, []string{}}
		for _, value := range header["value"].([]interface{}) {
			v := value.(bson.M)
			key := v["key"].(bson.ObjectId).Hex()
			val := v["val"].(string)
			hds[0] = append(hds[0], val)
			hds[1] = append(hds[1], key)
		 	selector[val + ".old"] = 0
		}
		return hds, selector, nil
	} else {
		return nil, nil, err
	}
}

func GetAllData() map[string]interface{} {
	db := S.DB("database")
	header, selector, err := GetDataHeaderAndSelector()
	var data []bson.M
	if err == nil {
		err = db.C("table").Find(bson.M{"deleted": false}).Select(selector).All(&data)
	}

	m := make(map[string]interface{})

	if err != nil {
		m["success"] = false
		m["reason"] = "数据库错误"
	} else {
		content, ids := convertData(data, header)
		m["header"] = header[0]
		m["content"] = content
		m["ids"] = ids
	}
	return m
}

func convertData(data []bson.M, headers [][]string) ([][]string, []string) {
	content := make([][]string, len(data))
	ids := make([]string, len(data))
	for line_no, value := range data {
		line := make([]string,len(headers[1]))
		for idx, header := range headers[1] {
			var v string
			if m, ok := value[header].(bson.M); ok {
				v = m["value"].(string)
			}
			line[idx] = v
		}
		content[line_no] = line
		ids[line_no] = value["_id"].(bson.ObjectId).Hex()
	}
	return content, ids
}
