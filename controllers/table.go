package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/models"
)

type TableController struct {
	beego.Controller
}

func (c *TableController) Get() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")   
	db := models.S.DB("database")
	var header interface{}
	err := db.C("column_name").Find(bson.M{}).Select(bson.M{"_id":0,"cur":1}).One(&header)
	var data []interface{}
	if err == nil {
		selector := make(bson.M)
		selector["_id"] = 0
		for _, value := range header.(bson.M)["cur"].([]interface{}) {
		 	selector[value.(string) + ".old"] = 0
		}
		err = db.C("data").Find(bson.M{}).Select(selector).All(&data)
	}
	m := make(map[string]interface{})

	if err != nil {
		m["success"] = false
		m["reason"] = err
	} else {
		var content [][]string
		for _, value := range data {
			var line []string
			for _, hd := range header.(bson.M)["cur"].([]interface{}) {
				line = append(line, value.(bson.M)[hd.(string)].(bson.M)["cur"].(string))
			}
			content = append(content, line)
		}
		m["header"] = header.(bson.M)["cur"]
		m["content"] = content
	}

	c.Data["json"] = m

}
