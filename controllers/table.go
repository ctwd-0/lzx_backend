package controllers

import (
	"github.com/astaxie/beego"
	"lzx_backend/models"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TableController struct {
	beego.Controller
}

func (c *TableController) InitTable() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")

	c.Data["json"] = models.GetAllData()
}

func updateColumn(new_header []string, author string) string{
	db := models.S.DB("database")

	var r bson.M
	iter := db.C("column").Find(nil).Iter()
	for iter.Next(&r) {
		err := db.C("column").UpdateId(r["_id"], bson.M{
			"$push":bson.M{"old":bson.M{"value":r["value"],"modified":r["modified"],"author":r["author"]}},
			"$set":bson.M{"value": new_header, "modified":time.Now(), "author": author},
		})
		if err != nil {
			return "更新表头失败"
		}
	}
	return ""
}

func removeColumn(column_name, author string) string {
	reason := ""
	header, _, err := models.GetDataHeaderAndSelector()

	if err != nil {
		reason = "获取表头失败"
	}

	new_header := []string{}
	if reason == "" {
		for _, value := range header {
			if value != column_name {
				new_header = append(new_header, value)
			}
		}
		if len(new_header) == len(header) {
			reason = "不存在的列"
		}
	}

	if reason == "" {
		reason = updateColumn(new_header, author)
	}

	return reason
}

func addColumn(column_name, author string) string {
	reason := ""
	header, _, err := models.GetDataHeaderAndSelector()

	if err != nil {
		reason = "获取表头失败"
	}

	if reason == "" {
		for _, val := range header {
			if val == column_name {
				reason = "重复的列"
				break
			}
		}
	}

	if reason == "" {
		reason = updateColumn(append(header, column_name), author)
	}

	return reason
}

func renameColumn(old_column_name, new_column_name, author string) string {
	reason := ""
	header, _, err := models.GetDataHeaderAndSelector()

	if err != nil {
		reason = "获取表头失败"
	}
	
	index := -1
	exist := false
	if reason == "" {
		for idx, val := range header {
			if val == old_column_name {
				index = idx
			} else if val == new_column_name {
				exist = true
			}
		}
		if index == -1 {
			reason = "未找到旧名字"
		} else if exist{ 
			reason = "重复的名字"
		} else {
			header[index] = new_column_name
		}
	}
	
	if reason == "" {
		reason = updateColumn(header, author)
	}

	return reason
}

func (c *TableController) RemoveColumn() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	column_name := c.GetString("column")
	author := "system_test"
	if column_name == "" {
		reason = "参数错误"
	}
	
	if reason == "" {
		reason = removeColumn(column_name, author)
	}

	if reason == "" {
		c.Data["json"] = bson.M{"success":true}
	} else {
		c.Data["json"] = bson.M{"success":false, "reason":reason}
	}
}

func (c *TableController) AddColumn() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	column_name := c.GetString("column")
	author := "system_test"
	if column_name == "" {
		reason = "参数错误"
	}
	
	if reason == "" {
		reason = addColumn(column_name, author)
	}

	if reason == "" {
		_, err := db.C("table").UpdateAll(bson.M{column_name:bson.M{"$exists": false}},bson.M{
			"$set":bson.M{column_name:bson.M{"value": "","author":"new_column", "modified":time.Now(), "old":[]bson.M{}}},
		})
		if err != nil {
			reason = "添加列失败"
		}
	}

	if reason == "" {
		c.Data["json"] = bson.M{"success":true}
	} else {
		c.Data["json"] = bson.M{"success":false, "reason":reason}
	}
}

func (c *TableController) RenameColumn() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	new_column_name := c.GetString("new")
	old_column_name := c.GetString("old")
	author := "system_test"
	if new_column_name == "" || old_column_name == "" || old_column_name == new_column_name {
		reason = "参数错误"
	}
	
	if reason == "" {
		reason = renameColumn(old_column_name, new_column_name, author)
	}

	if reason == "" {
		var r bson.M
		iter := db.C("table").Find(bson.M{old_column_name:bson.M{"$exists": true}}).Select(bson.M{old_column_name:1}).Iter()
		for iter.Next(&r) {
			err := db.C("table").UpdateId(r["_id"],bson.M{
				"$set":bson.M{new_column_name:r[old_column_name]},
				"$unset":bson.M{old_column_name:""},
			})
			if err != nil {
				reason = "更新列名失败(数据库可能不一致)"
			}
		}

	}

	if reason == "" {
		c.Data["json"] = bson.M{"success":true}
	} else {
		c.Data["json"] = bson.M{"success":false, "reason":reason}
	}
}

func (c *TableController) UpdateValue() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	id_hex := c.GetString("id")
	column_name := c.GetString("column")
	value := c.GetString("value")
	author := "system_test"
	reason := ""
	if column_name == "" || !bson.IsObjectIdHex(id_hex) {
		reason = "参数错误"
	}

	if reason == "" {
		header, _, err := models.GetDataHeaderAndSelector()
		if err != nil {
			reason = "获取表头失败"
		} else {
			reason = "不存在的列"
			for _, val := range header {
				if val == column_name {
					reason = ""
					break
				}
			}
		}
	}

	if reason == "" {
		var r bson.M
		err := db.C("table").FindId(bson.ObjectIdHex(id_hex)).One(&r)
		if err != nil {
			reason = "查找失败"
		} else {
			err = db.C("table").UpdateId(r["_id"], bson.M{
				"$push":bson.M{column_name + ".old": bson.M{
					"value":r[column_name].(bson.M)["value"],
					"modified":r[column_name].(bson.M)["modified"],
					"author":r[column_name].(bson.M)["author"],
				}},
				"$set":bson.M{
					column_name+".value": value,
					column_name+".modified":time.Now(),
					column_name+".author": author,
				},
			})
			if err != nil {
				reason = "更新失败"
			}
		}
	}
	if reason == "" {
		c.Data["json"] = bson.M{"success":true}
	} else {
		c.Data["json"] = bson.M{"success":false, "reason":reason}
	}
}