package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/utils"
	"lzx_backend/models"
	"net/http"
)

const admin_name = "admin"

type AdminController struct {
	beego.Controller
}

//管理员获取所有的用户
func allUsers() ([]bson.M, string) {
	db := models.S.DB("database")
	var users []bson.M
	err := db.C("user").Pipe([]bson.M{
		{"$match":bson.M{"deleted":false, "name":bson.M{"$ne":admin_name}}},
		{"$project":bson.M{"_id":0,"name":1,"write":1, "banned":1}},
	}).All(&users)
	if err != nil {
		return []bson.M{}, "数据库错误"
	} else {
		return users, ""
	}
}

//管理员获取所有的用户。
func (c *AdminController) InitUser() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	reason := ""
	
	if c.GetSession("name") != admin_name || c.GetSession("admin") != true {
		reason = "权限不足"
	}

	var users []bson.M
	if reason == "" {
		users, reason = allUsers()
	}
	
	m := SimpleReturn(reason)
	if reason == "" {
		m["users"] = users;
	}

	c.Data["json"] = m
}

//管理员增加用户
func (c *AdminController) AddUser() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	db := models.S.DB("database")
	reason := ""
	name := c.GetString("name")
	pass_md5 := c.GetString("password")
	
	if c.GetSession("name") != admin_name || c.GetSession("admin") != true {
		reason = "权限不足"
	}

	if reason == "" {
		if name == "" || !utils.IsLowerMD5(pass_md5) {
			reason = "参数错误"
		}
	}

	if reason == "" {
		reason = userNotExist(name)
	}

	if reason == "" {
		err := db.C("user").Insert(bson.M{
			"name":name, "password": utils.SaltPassword(pass_md5),
			"write": true, "banned": false, "deleted": false})
		if err != nil {
			reason = "数据库错误"
		}
	}
	
	var users []bson.M
	if reason == "" {
		users, reason = allUsers()
	}

	c.Data["json"] = SimpleReturn(reason)
	if reason == "" {
		c.Data["json"].(map[string]interface{})["users"] = users
	}
}

//管理员删除用户
func (c *AdminController) DeleteUser() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	db := models.S.DB("database")
	reason := ""
	name := c.GetString("name")
	
	if c.GetSession("name") != admin_name || c.GetSession("admin") != true {
		reason = "权限不足"
	}

	if reason == "" {
		if name == "" || name == admin_name {
			reason = "参数错误"
		}
	}

	if reason == "" {
		reason = userExist(name)
	}

	if reason == "" {
		err := db.C("user").Update(bson.M{"name":name, "deleted":false},
			bson.M{"$set":bson.M{"deleted":true}})
			if err != nil {
			reason = "数据库错误"
		}
	}

	var users []bson.M
	if reason == "" {
		users, reason = allUsers()
	}

	c.Data["json"] = SimpleReturn(reason)
	if reason == "" {
		c.Data["json"].(map[string]interface{})["users"] = users
	}
}

//管理员重置密码
func (c *AdminController) ChangePassword() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	db := models.S.DB("database")
	reason := ""
	name := c.GetString("name")
	pass_md5 := c.GetString("password")

	if c.GetSession("name") != admin_name || c.GetSession("admin") != true {
		reason = "权限不足"
	}

	if reason == "" {
		if name == "" || !utils.IsLowerMD5(pass_md5) {
			reason = "参数错误"
		}
	}

	if reason == "" {
		reason = userExist(name)
	}

	if reason == "" {
		err := db.C("user").Update(bson.M{"name":name, "deleted":false},
			bson.M{"$set":bson.M{"password": utils.SaltPassword(pass_md5)}})
			if err != nil {
			reason = "数据库错误"
		}
	}

	c.Data["json"] = SimpleReturn(reason)
}

//管理员更新用户信息
func (c *AdminController) UpdateUser() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	db := models.S.DB("database")
	reason := ""
	name := c.GetString("name")
	banned, err_b := c.GetBool("banned")
	write, err_w := c.GetBool("write")

	if c.GetSession("name") != admin_name || c.GetSession("admin") != true {
		reason = "权限不足"
	}
	
	if reason == "" {
		if name == "" || err_b != nil || err_w != nil {
			reason = "参数错误"
		}
	}

	if reason == "" {
		reason = userExist(name)
	}

	if reason == "" {
		err := db.C("user").Update(bson.M{"name":name, "deleted":false},
			bson.M{"$set":bson.M{"write":write, "banned": banned}})
			if err != nil {
			reason = "数据库错误"
		}
	}

	var users []bson.M
	if reason == "" {
		users, reason = allUsers()
	}

	c.Data["json"] = SimpleReturn(reason)
	if reason == "" {
		c.Data["json"].(map[string]interface{})["users"] = users
	}
}

//管理员页面
func (c *AdminController) Admin() {
	if c.GetSession("name") == admin_name {
		http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, "dist/admin.html")
	} else{
		c.Ctx.Redirect(302, "/")
	}
}
