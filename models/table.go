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
		db.C("column").Insert(bson.M{
			"value": headers,
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
					m[headers[index]] = bson.M{
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
		m["reason"] = "db error"
	} else {
		content := convertData(data, header)
		m["header"] = header
		m["content"] = content
	}
	return m
}

func GetDataHeaderAndSelector() ([]string, bson.M, error) {
	db := S.DB("database")
	var header bson.M
	err := db.C("column").Find(bson.M{}).Select(bson.M{"_id":0,"value":1}).One(&header)
	if err == nil {
		selector := bson.M{}
		selector["_id"] = 0
		h := []string{}
		for _, value := range header["value"].([]interface{}) {
			h = append(h, value.(string) )
		 	selector[value.(string) + ".old"] = 0
		}
		return h, selector, nil
	} else {
		return nil, nil, err
	}
}

func GetAllData() map[string]interface{} {
	db := S.DB("database")
	header, selector, err := GetDataHeaderAndSelector()
	var data []bson.M
	fmt.Println(time.Now())
	if err == nil {
		err = db.C("table").Find(bson.M{"deleted": false}).Select(selector).All(&data)
	}
	fmt.Println(time.Now())

	m := make(map[string]interface{})

	if err != nil {
		m["success"] = false
		m["reason"] = "db error"
	} else {
		fmt.Println(time.Now())
		content := convertData(data, header)
		fmt.Println(time.Now())
		m["header"] = header
		m["content"] = content
	}
	return m
}

func convertData(data []bson.M, headers []string) [][]string {
	content := make([][]string, 0)
	for _, value := range data {
		var line []string
		for _, header := range headers {
			var v string
			if m, ok := value[header].(bson.M); ok {
				v = m["value"].(string)
			}
			line = append(line, v)
		}
		content = append(content, line)
	}
	return content
}
