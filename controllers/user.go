package controllers

import (
	//"fmt"
	"lzx_backend/models"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/utils"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

//返回用户存在的原因
func userNotExist(name string) string {
	db := models.S.DB("database")
	count, err := db.C("user").Find(bson.M{"name":name, "deleted": false}).Count()
	if err != nil {
		return "数据库错误"
	} else if count != 0 {
		return "名称重复"
	}
	return ""
}

//返回用户不存在的原因
func userExist(name string) string {
	db := models.S.DB("database")
	count, err := db.C("user").Find(bson.M{"name":name, "deleted": false}).Count()
	if err != nil {
		return "数据库错误"
	} else if count == 0 {
		return "未找到"
	}
	return ""
}

//当只需返回操作是否成功时，使用该函数构造返回值
func SimpleReturn(reason string) map[string]interface{} {
	m := map[string]interface{}{}
	if reason == "" {
		m["success"] = true
	} else {
		m["success"] = false
		m["reason"] = reason
	}
	return m
}

//使用旧密码重设密码。未使用。
func (c *UserController) ChangePassword() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	db := models.S.DB("database")
	reason := ""
	name := c.GetString("name")
	old_md5 := c.GetString("old")
	new_md5 := c.GetString("new")

	if name == "" || !utils.IsLowerMD5(old_md5) || !utils.IsLowerMD5(new_md5) {
		reason = "参数错误"
	}

	if reason == "" {
		var user bson.M
		err := db.C("user").Find(bson.M{"name":name, "deleted": false}).One(&user)
		if err != nil {
			reason = "数据库错误"
		} else if user == nil {
			reason = "不存在的用户"
		} else if user["password"] != utils.SaltPassword(old_md5) {
			reason = "旧密码错误"
		}
	}

	if reason == "" {
		err := db.C("user").Update(bson.M{"name":name, "deleted":false},
			bson.M{"$set":bson.M{"password":new_md5}})
			if err != nil {
			reason = "数据库错误"
		}
	}

	c.Data["json"] = SimpleReturn(reason)
}

//用户登陆，设置session
func (c *UserController) Login() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	
	db := models.S.DB("database")
	reason := ""
	name := c.GetString("name")
	pass_md5 := c.GetString("password")

	if name == "" || !utils.IsLowerMD5(pass_md5) {
		reason = "参数错误"
	}

	var user bson.M
	if reason == "" {
		err := db.C("user").Find(bson.M{"name":name, "deleted": false}).One(&user)
		if err != nil {
			reason = err.Error()//"数据库错误"
		} else if user == nil {
			reason = "不存在的用户"
		} else if user["password"] != utils.SaltPassword(pass_md5) {
			reason = "密码错误"
		} else if user["banned"] == true {
			reason = "用户被管理员禁止"
		}
	}

	if reason == "" {
		c.SetSession("name", user["name"])
		c.SetSession("write", user["write"])
		if user["admin"]  == true {
			c.SetSession("admin", true)
		}
	}

	c.Data["json"] = SimpleReturn(reason)
}

//用户登出
func (c *UserController) Logout() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

	c.DelSession("name")
	c.DelSession("write")
	c.DelSession("admin")

	c.Data["json"] = bson.M{
		"success": true,
	}
}

//检查用户是否有写权限
func CheckWrite(write, author interface{}) string {
	if write == nil || author == nil {
		return "未登陆无法操作。"
	} else if write.(bool) == false {
		return "无写权限无法操作。"
	} else { 
		return ""
	}
}