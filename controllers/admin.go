package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/utils"
	"lzx_backend/models"
)

const admin_name = "admin"

type AdminController struct {
	beego.Controller
}

func (c *AdminController) InitUser() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	
	db := models.S.DB("database")
	reason := ""
	

	if c.GetSession("name") != admin_name || c.GetSession("admin") != true {
		reason = "权限不足"
	}

	var users []bson.M
	if reason == "" {
		err := db.C("user").Pipe([]bson.M{
			{"$match":bson.M{"deleted":false, "name":bson.M{"$ne":admin_name}}},
			{"$project":bson.M{"_id":0,"name":1,"write":1, "banned":1}},
		}).All(&users)
		if err != nil {
			reason = "数据库错误"
		}
	}
	
	m := SimpleReturn(reason)
	if reason == "" {
		m["users"] = users;
	}

	c.Data["json"] = m
}

func (c *AdminController) AddUser() {
	defer c.ServeJSON()
	//c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	db := models.S.DB("database")
	reason := ""
	name := c.GetString("name")
	pass_md5 := c.GetString("password")
	
	if c.GetSession("name") != admin_name || c.GetSession("admin") != true {
		reason = "权限不足"
	}

	if name == "" || !utils.IsLowerMD5(pass_md5) {
		reason = "参数错误"
	}

	if reason == "" {
		reason = userExist(name)
	}

	if reason == "" {
		err := db.C("user").Insert(bson.M{
			"name":name, "password": utils.SaltPassword(pass_md5),
			"write": true, "banned": false, "deleted": false})
		if err != nil {
			reason = "数据库错误"
		}
	}

	c.Data["json"] = SimpleReturn(reason)
}