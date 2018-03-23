package controllers

import (
	//"fmt"
	"encoding/json"
	"lzx_backend/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
)

type FilterController struct {
	beego.Controller
}

func pureHeaders() ([]string, error) {
	hds, _, err := models.GetDataHeaderAndSelector()
	if err != nil {
		return []string{}, err
	} else {
		header := []string{}
		for _, val := range hds.([]interface{}) {
			if val == "构件编号（表单中显示）" {
				header = append(header, "构件编号")
			} else if val != "模型编号（rhino中对应编号，表单中表头、值均不显示）"{
				header = append(header, val.(string))
			}
		}
		return header, nil
	}
}

func (c *FilterController) InitFilter() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	header, err := pureHeaders()
	if err != nil {
		reason = "数据库错误"
	}

	if reason == "" {
		count, err := db.C("filter").Find(bson.M{}).Count()
		if err != nil {
			reason = "数据库错误"
		} else if count == 0 {
			addFilter("默认", []string{"构件编号"}, true)
		}
	}

	var result []bson.M
	if reason == "" {
		err = db.C("filter").Pipe([]bson.M{
			{"$match": bson.M{"deleted": false}},
			{"$project": bson.M{"_id":0,"name":1,"model":1,"defualt":1}},
		}).All(&result)
		if err != nil || result == nil {
			reason = "数据库错误"
		}
	}

	m := map[string]interface{}{}
	if reason == "" {
		m["header"] = header
		m["filter"] = result
	} else {
		m["success"] = false
		m["reason"] = reason
	}
	c.Data["json"] = m
}	

func addFilter(name string, model []string, defualt bool) error {
	db := models.S.DB("database")
	return db.C("filter").Insert(bson.M{"name":name, "model": model, "defualt":defualt, "deleted": false})
}

func (c *FilterController) AddFilter() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	//reason := ""
	//header, err := pureHeaders()
}

func (c *FilterController) DeleteFilter() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
}

func (c *FilterController) UpdateFilter() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	name := c.GetString("name")
	model_str :=c.GetString("model")
	if name == "" || model_str == "" {
		reason = "参数不能为空"
	}
	var model_i interface{}
	var model []string
	if reason == "" {
		err := json.Unmarshal([]byte(model_str), &model_i)
		if err != nil {
			reason = "格式错误"
		}
		if model_i, ok := model_i.([]interface{}); ok {
			for _, val := range model_i {
				if val, ok := val.(string); ok {
					model = append(model, val)
				}
			}
		} else {
			reason = "格式错误"
		}
		if reason == "" {
			if model == nil || len(model) == 0 {
				reason = "格式错误"
			}
		}
	}

	if reason == "" {
		_, err := db.C("filter").UpdateAll(bson.M{},bson.M{"$set":bson.M{"defualt":false}})
		if err != nil {
			reason = "数据库错误"
		} 
	}

	if reason == "" {
		err := db.C("filter").Update(bson.M{"name":name}, bson.M{"$set":bson.M{"defualt":true, "model":model}})
		if err != nil {
			reason = "数据库错误"
		}
	}

	m := map[string]interface{}{}
	if reason == "" {
		m["success"] = true
	} else {
		m["success"] = false
		m["reason"] = reason
	}
	c.Data["json"] = m
}