package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/models"
	"github.com/satori/go.uuid"
)

type FileController struct {
	beego.Controller
}

func (c *FileController) Upload() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	model_id := c.GetString("model_id")
	category := c.GetString("category")

	if model_id == "" || category == "" {
		c.Data["json"] = bson.M{"success":false, "reason": "参数错误"}
	} else {
		file, header, _ := c.GetFile("file")

		filename := header.Filename

		uuid, _ := uuid.NewV4()

		go models.ProcessUploadedFile(file, filename, model_id, category, uuid)

		c.Data["json"] = bson.M{"success":true, "token": uuid}
	}
}

func (c *FileController) Update() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	id_hex := c.GetString("id")
	description := c.GetString("description")

	if description == "" || !bson.IsObjectIdHex(id_hex) {
		reason = "参数错误"
	}

	if reason == "" {
		err := db.C("file").UpdateId(bson.ObjectIdHex(id_hex), bson.M{
			"$set":bson.M{"description":description},
		})
		if err != nil {
			reason = "数据库错误"
		}
	}

	if reason == "" {
		c.Data["json"] = bson.M{"success":true}
	} else {
		c.Data["json"] = bson.M{"success":false, "reason": "参数错误"}
	}
}

func (c *FileController) Options() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	c.Data["json"] = bson.M{"success":true}
}

func allFiles(model_id, category string) ([]bson.M, string){
	db := models.S.DB("database")
	category_id, reason := models.ConvertName(model_id, category)

	data:= []bson.M{}
	if reason == "" {
		err := db.C("file").Find(bson.M{
			"model_id": model_id, "category": category_id, "deleted": false,
		}).Select(bson.M{
			"deleted":0,"created":0,"uuid":0,
			"original_md5":0,"thumbnail_md5":0,"thumbnail_saved_as":0,"original_saved_as":0,
		}).Sort("-created").All(&data)
		if err != nil {
			reason = "数据库错误"
		}
	}
	return data, reason
}

func (c *FileController) GetAll() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")   
	reason := ""
	model_id := c.GetString("model_id")
	category := c.GetString("category")

	if model_id == "" || category == "" {
		reason = "参数错误"
	}

	var data []bson.M
	if reason == "" {
		data, reason = allFiles(model_id, category)
	}

	m := map[string]interface{}{}
	if reason == "" {
		m["success"] = true
		m["files"] = data
	} else {
		m["success"] = false
		m["reason"] = reason
	}

	c.Data["json"] = m
}

func (c *FileController) Ready() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")   
	db := models.S.DB("database")
	reason := ""
	token := c.GetString("token")
	model_id := c.GetString("model_id")
	category := c.GetString("category")
	if token == ""  || model_id == "" || category == "" {
		reason = "参数错误"
	}

	if reason == "" {
		count, err := db.C("file").Find(bson.M{"uuid":token}).Count()
		if err != nil {
			reason = "数据库错误"
		} else if count == 0 {
			reason = "尚未准备好"
		}
	}

	var files []bson.M
	if reason == "" {
		files, reason = allFiles(model_id, category)
	}

	if reason == "" {
		c.Data["json"] = bson.M{"success":true, "files": files}
	} else {
		c.Data["json"] = bson.M{"success":false, "reason": reason}
	}
}
