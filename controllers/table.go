package controllers

import (
	"github.com/astaxie/beego"
	"lzx_backend/models"
)

type TableController struct {
	beego.Controller
}

func (c *TableController) Get() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")

	c.Data["json"] = models.GetAllData()
}
