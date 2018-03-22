package controllers

import (
	//"fmt"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/models"
	"github.com/astaxie/beego"
)

type QueryController struct {
	beego.Controller
}

// func (c *QueryController) Get() {
// 	defer c.ServeJSON()
// 	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")

// 	query_string := c.GetString("query")
// 	//fmt.Println(query_string)

// 	c.Data["json"] = models.QueryDataWithString(query_string)
// }

func getNames() map[string]interface{} {
	var result bson.M
	db := models.S.DB("database")
	err := db.C("query").Pipe([]bson.M{
		{"$group":bson.M{"_id":nil, "names":bson.M{"$push":"$name"}}},
	}).One(&result)
	m := make(map[string]interface{})
	if err != nil || result == nil {
		m["success"] = false
		m["reason"] = "数据库错误"
		m["names"] = []string{}
	} else {
		m["names"] = result["names"]
	}
	return m
}

func (c *QueryController) InitQuery() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")

	c.Data["json"] = getNames()
}

func (c *QueryController) SaveQuery() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	db := models.S.DB("database")
	name := c.GetString("name")
	query := c.GetString("query")
	reason := ""
	if name == "" || query == "" {
		reason = "参数错误"
	}
	var query_map interface{}
	err := json.Unmarshal([]byte(query), &query_map)
	if err != nil {
		reason = "查询格式错误"
	}

	if reason == "" {
		var result bson.M
		err := db.C("query").Find(bson.M{"name":name, "deleted": false}).One(&result)
		if err == nil && result != nil  {
			reason = "名称重复"
		}
	}

	if reason == "" {
		err = db.C("query").Insert(bson.M{"name":name, "query": query, "deleted": false})
		if err != nil {
			reason = "数据库错误"
		}
	}
	if reason == "" {
		m := getNames()
		if m["success"] == nil {
			m["success"] = true
		}
		c.Data["json"] = m
	} else {
		m := map[string]interface{}{}
		m["success"] = false
		m["reason"] = reason
		c.Data["json"] = m
	}
}