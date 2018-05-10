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

//上传文件
func (c *FileController) Upload() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	model_id := c.GetString("model_id")
	category := c.GetString("category")
	author := c.GetSession("name")
	
	reason := CheckWrite(c.GetSession("write"), author)
	if reason != ""  {
		c.Data["json"] = bson.M{"success":false, "reason": reason}
	} else if model_id == "" || category == "" {
		c.Data["json"] = bson.M{"success":false, "reason": "参数错误"}
	} else {
		file, header, _ := c.GetFile("file")

		filename := header.Filename

		uuid, _ := uuid.NewV4()

		go models.ProcessUploadedFile(file, filename, model_id, category, uuid)

		c.Data["json"] = bson.M{"success":true, "token": uuid}
	}
}

//更新文件描述
func (c *FileController) Update() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	id_hex := c.GetString("id")
	description := c.GetString("description")
	author := c.GetSession("name")
	
	reason := CheckWrite(c.GetSession("write"), author)
	if reason == "" {
		if description == "" || !bson.IsObjectIdHex(id_hex) {
			reason = "参数错误"
		}
	}

	if reason == "" {
		err := db.C("file").UpdateId(bson.ObjectIdHex(id_hex), bson.M{
			"$set":bson.M{"description":description, "author": author},
		})
		if err != nil {
			reason = "数据库错误"
		}
	}

	c.Data["json"] = SimpleReturn(reason)
}

func (c *FileController) Options() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	c.Data["json"] = bson.M{"success":true}
}

//删除文件
func (c *FileController) Remove() {
	defer c.ServeJSON()
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	id_hex := c.GetString("id")
	author := c.GetSession("name")
	
	reason := CheckWrite(c.GetSession("write"), author)
	if reason == "" {
		if !bson.IsObjectIdHex(id_hex) {
			reason = "参数错误"
		}
	}

	if reason == "" {
		err := db.C("file").UpdateId(bson.ObjectIdHex(id_hex), bson.M{
			"$set":bson.M{"deleted":true, "deleted_by": author},
		})
		if err != nil {
			reason = "数据库错误"
		}
	}

	c.Data["json"] = SimpleReturn(reason)
}

//下载文件
func (c *FileController) Download() {
	db := models.S.DB("database")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	id_hex := c.GetString("id")
	author := c.GetSession("name")

	reason := CheckWrite(c.GetSession("write"), author)
	if reason == "" {
		if !bson.IsObjectIdHex(id_hex) {
			reason = "参数错误"
		}
	}

	var data bson.M
	if reason == "" {
		err := db.C("file").FindId(bson.ObjectIdHex(id_hex)).One(&data)
		if err != nil {
			reason = "数据库错误"
		}
	}

	if reason == "" {
		path := data["original_path"].(string)
		name := data["original_name"].(string)
		if path[0] == '/' {
			path = path[1:]
		}
		c.Ctx.Output.Download(path, name)
	} else {
		writer := c.Ctx.ResponseWriter
		writer.Header().Set("Content-Disposition", "attachment; filename=" + reason + ".txt")
		writer.Header().Set("Content-Description", "File Transfer")
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Header().Set("Content-Transfer-Encoding", "binary")
		writer.Header().Set("Expires", "0")
		writer.Header().Set("Cache-Control", "must-revalidate")
		writer.Header().Set("Pragma", "public")
		c.Ctx.ResponseWriter.Write([]byte(reason))
	}
}

//获取某个构件某个文件夹下的所有文件
func (c *FileController) GetAll() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	reason := ""
	model_id := c.GetString("model_id")
	category := c.GetString("category")

	if model_id == "" || category == "" {
		reason = "参数错误"
	}

	var files []bson.M
	if reason == "" {
		files, reason = models.AllFiles(model_id, category)
	}

	c.Data["json"] = SimpleReturn(reason)
	if reason == "" {
		c.Data["json"].(map[string]interface{})["files"] = files
	}
}

//测试某个上传的文件是否已已经就绪
func (c *FileController) Ready() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	db := models.S.DB("database")
	reason := ""
	token := c.GetString("token")
	model_id := c.GetString("model_id")
	category := c.GetString("category")

	if token == ""  || model_id == "" || category == "" {
		reason = "参数错误"
	}

	var data bson.M
	if reason == "" {
		err := db.C("file").Find(bson.M{"uuid":token}).One(&data)
		if err != nil || data == nil {
			reason = "尚未准备好"
		}
	}

	var files []bson.M
	if reason == "" {
		files, reason = models.AllFiles(model_id, category)
	}

	c.Data["json"] = SimpleReturn(reason)
	if reason == "" {
		c.Data["json"].(map[string]interface{})["files"] = files
		if data["reason"] != nil {
			c.Data["json"].(map[string]interface{})["finish_with_error"] = true
			c.Data["json"].(map[string]interface{})["finish_with_error_reason"] = data["reason"]
		}
	}
}
