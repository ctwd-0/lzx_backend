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


//初始化，返回数据表内容。
func (c *TableController) InitTable() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

	c.Data["json"] = models.GetAllData()
}

//数据操作，更新列名
func updateColumn(new_header [][]string, author string) string {
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

//数据操作，删除列
func removeColumn(column_name, author string) string {
	reason := ""
	header, _, err := models.GetDataHeaderAndSelector()

	if err != nil {
		reason = "获取表头失败"
	}

	new_header := [][]string{[]string{},[]string{}}
	if reason == "" {
		if len(header[0]) == 1 {
			return "不能删除最后一列"
		}
		for idx, value := range header[0] {
			if value != column_name {
				new_header[0] = append(new_header[0], value)
				new_header[1] = append(new_header[1], header[1][idx])
			}
		}
		if len(new_header[0]) == len(header[0]) {
			reason = "不存在的列名称"
		}
	}

	if reason == "" {
		reason = updateColumn(new_header, author)
	}

	return reason
}

//数据操作，增加列
func addColumn(column_name, author string) (string, string) {
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

	var column_id string
	if reason == "" {
		column_id = bson.NewObjectId().Hex()
		header[0] = append(header[0], column_name)
		header[1] = append(header[1], column_id)
		reason = updateColumn(header, author)
	}

	return column_id, reason
}

//数据操作，重命名列
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

//生成返回数据。包含数据表的表头
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

//删除列名
func (c *TableController) RemoveColumn() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	column_name := c.GetString("column")
	author := c.GetSession("name")

	reason := CheckWrite(c.GetSession("write"), author)
	if reason == "" && column_name == "" {
		reason = "参数错误"
	}

	if reason == "" {
		if column_name == "构件编号" {
			reason = "不能删除构件编号列"
		}
	}
	
	if reason == "" {
		reason = removeColumn(column_name, author.(string))
	}

	c.Data["json"] = returnHeader(reason)
}

//增加列名
func (c *TableController) AddColumn() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	column_name := c.GetString("column")
	author := c.GetSession("name")

	reason := CheckWrite(c.GetSession("write"), author)
	if reason == "" && column_name == "" {
		reason = "参数错误"
	}
	
	var column_id string
	if reason == "" {
		column_id, reason = addColumn(column_name, author.(string))
	}

	if reason == "" {
		_, err := db.C("table").UpdateAll(bson.M{column_id:bson.M{"$exists": false}},bson.M{
			"$set":bson.M{column_id:bson.M{"value": "","author":"new_column", "modified":time.Now(), "old":[]bson.M{}}},
		})
		if err != nil {
			reason = "添加列失败"
		}
	}

	c.Data["json"] = returnHeader(reason)
}

//重命名列
func (c *TableController) RenameColumn() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	new_column_name := c.GetString("new")
	old_column_name := c.GetString("old")
	author := c.GetSession("name")
	
	reason := CheckWrite(c.GetSession("write"), author)
	if reason == "" {
		if new_column_name == "" || old_column_name == "" || old_column_name == new_column_name {
			reason = "参数错误"
		}
	}
	
	if reason == "" {
		if old_column_name == "构件编号" {
			reason = "不能重命名构件编号列"
		}
	}

	if reason == "" {
		reason = renameColumn(old_column_name, new_column_name, author.(string))
	}

	c.Data["json"] = returnHeader(reason)
}

//更新表中数值
func (c *TableController) UpdateValue() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	id_hex := c.GetString("id")
	column_name := c.GetString("column")
	value := c.GetString("value")
	author := c.GetSession("name")
	
	reason := CheckWrite(c.GetSession("write"), author)
	if reason == "" {
		if column_name == "" || !bson.IsObjectIdHex(id_hex) {
			reason = "参数错误"
		}
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
