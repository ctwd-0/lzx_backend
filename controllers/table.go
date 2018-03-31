package controllers

import (
	"github.com/astaxie/beego"
	"lzx_backend/models"
	"lzx_backend/utils"
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

func updateColumn(new_header [][]string, author string) string{
	db := models.S.DB("database")
	data := utils.StringToData(new_header)
	var r bson.M
	iter := db.C("column").Find(nil).Iter()
	for iter.Next(&r) {
		err := db.C("column").UpdateId(r["_id"], bson.M{
			"$push":bson.M{"old":bson.M{"value":r["value"],"modified":r["modified"],"author":r["author"]}},
			"$set":bson.M{"value": data, "modified":time.Now(), "author": author},
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

	new_header := [][]string{[]string{},[]string{}}
	if reason == "" {
		for idx, value := range header[0] {
			if value != column_name {
				new_header[0] = append(new_header[0], value)
				new_header[1] = append(new_header[1], header[1][idx])
			}
		}
		if len(new_header) == len(header) {
			reason = "不存在的列名称"
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
		for _, val := range header[0] {
			if val == column_name {
				reason = "重复的列"
				break
			}
		}
	}

	if reason == "" {
		header[0] = append(header[0], column_name)
		header[1] = append(header[1], bson.NewObjectId().Hex())
		reason = updateColumn(header, author)
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
		for idx, val := range header[0] {
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
			header[0][index] = new_column_name
		}
	}
	
	if reason == "" {
		reason = updateColumn(header, author)
	}

	return reason
}

func returnHeader(reason string) bson.M {
	header := []string{}
	if reason == "" {
		header, reason = models.GetDataHeader()
	}

	if reason == "" {
		return bson.M{"success":true, "header": header}
	} else {
		return bson.M{"success":false, "reason": reason}
	}
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

	c.Data["json"] = returnHeader(reason)
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

	c.Data["json"] = returnHeader(reason)
}

func (c *TableController) RenameColumn() {
	defer c.ServeJSON()
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

	c.Data["json"] = returnHeader(reason)
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

	var column_id string
	if reason == "" {
		header, _, err := models.GetDataHeaderAndSelector()
		if err != nil {
			reason = "获取表头失败"
		} else {
			reason = "不存在的列"
			for idx, val := range header[0] {
				if val == column_name {
					reason = ""
					column_id = header[1][idx]
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
				"$push":bson.M{column_id + ".old": bson.M{
					"value":r[column_id].(bson.M)["value"],
					"modified":r[column_id].(bson.M)["modified"],
					"author":r[column_id].(bson.M)["author"],
				}},
				"$set":bson.M{
					column_id+".value": value,
					column_id+".modified":time.Now(),
					column_id+".author": author,
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
