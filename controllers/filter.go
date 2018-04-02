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

func pureHeaders() ([]string, string) {
	hds, reason := models.GetDataHeader()
	if reason != "" {
		return []string{}, reason
	} else {
		header := []string{}
		for _, val := range hds {
			if val == "构件编号（表单中显示）" {
				header = append(header, "构件编号")
			} else if val != "模型编号（rhino中对应编号，表单中表头、值均不显示）" && val != "模型编号" {
				header = append(header, val)
			}
		}
		return header, ""
	}
}

func initData() map[string]interface{} {
	db := models.S.DB("database")

	header, reason := pureHeaders()

	if reason == "" {
		count, err := db.C("filter").Find(bson.M{"deleted": false}).Count()
		if err != nil {
			reason = "数据库错误"
		} else if count == 0 {
			addFilter("默认", []string{"构件编号"}, true)
		}
	}

	var result []bson.M
	if reason == "" {
		err := db.C("filter").Pipe([]bson.M{
			{"$match": bson.M{"deleted": false}},
			{"$project": bson.M{"_id":0,"name":1,"model":1,"default":1}},
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

	return m
}

func getStringArray(input string) ([]string, string){
	reason := ""
	var output_i interface{}
	output:= []string{}
	err := json.Unmarshal([]byte(input), &output_i)
	if err != nil {
		reason = "格式错误"
	}
	if output_i, ok := output_i.([]interface{}); ok {
		for _, val := range output_i {
			if val, ok := val.(string); ok {
				output = append(output, val)
			}
		}
	} else {
		reason = "格式错误"
	}
	if reason == "" {
		if output == nil || len(output) == 0 {
			reason = "格式错误"
		}
	}
	return output, reason
}

func addFilter(name string, model []string, dft bool) error {
	db := models.S.DB("database")
	return db.C("filter").Insert(bson.M{"name":name, "model": model, "default":dft, "deleted": false})
}

func (c *FilterController) InitFilter() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	
	c.Data["json"] = initData()
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

	var model []string
	if reason == "" {
		model, reason = getStringArray(model_str)
	}

	if reason == "" {
		_, err := db.C("filter").UpdateAll(bson.M{"deleted":false},
			bson.M{"$set":bson.M{"default":false}})
		if err != nil {
			reason = "数据库错误"
		} 
	}

	if reason == "" {
		err := db.C("filter").Update(bson.M{"name":name, "deleted": false},
			bson.M{"$set":bson.M{"default":true, "model":model}})
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

func (c *FilterController) DeleteFilter() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	db := models.S.DB("database")
	reason := ""

	name := c.GetString("name")
	if name == "" {
		reason = "参数不能为空"
	}

	if reason == "" {
		count, err := db.C("filter").Find(bson.M{"deleted": false}).Count()
		if err != nil {
			reason = "数据库错误"
		} else if count == 1 {
			reason = "不能删除最后一个表头筛选条件"
		}
	}

	if reason == "" {
		err := db.C("filter").Update(bson.M{"name":name, "deleted":false}, bson.M{"$set":bson.M{"deleted":true}})
		if err != nil {
			reason = "数据库错误"
		}
	}

	if reason == "" {
		count, err := db.C("filter").Find(bson.M{"deleted": false, "default":true}).Count()
		if err != nil {
			reason = "数据库错误"
		} else if count == 0 {
			err := db.C("filter").Update(bson.M{"deleted":false}, bson.M{"$set":bson.M{"default":true}})
			if err != nil {
				reason = "数据库错误"
			}
		}
	}

	m := map[string]interface{}{}
	if reason == "" {
		m = initData()
		m["success"] = true
	} else {
		m["success"] = false
		m["reason"] = reason
	}

	c.Data["json"] = m
}

func (c *FilterController) AddFilter() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	db := models.S.DB("database")
	reason := ""
	name := c.GetString("name")
	model_str :=c.GetString("model")
	if name == "" || model_str == "" {
		reason = "参数不能为空"
	}

	var model []string
	if reason == "" {
		model, reason = getStringArray(model_str)
	}

	if reason == "" {
		count, err := db.C("filter").Find(bson.M{"name":name, "deleted":false}).Count()
		if err != nil {
			reason = "数据库错误"
		} else if count != 0 {
			reason = "名称重复"
		}
	}

	if reason == "" {
		_, err := db.C("filter").UpdateAll(bson.M{"deleted":false},bson.M{"$set":bson.M{"default":false}})
		if err != nil {
			reason = "数据库错误"
		} 
	}

	if reason == "" {
		err := addFilter(name, model, true)
		if err != nil {
			reason = "数据库错误"
		}
	}

	m := map[string]interface{}{}
	if reason == "" {
		m = initData()
		m["success"] = true
	} else {
		m["success"] = false
		m["reason"] = reason
	}

	c.Data["json"] = m
}
