package controllers

import (
	//"fmt"
	"lzx_backend/models"
	"github.com/astaxie/beego"
)

type SearchController struct {
	beego.Controller
}

func (c *SearchController) Search() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	
	query_string := c.GetString("query")

	c.Data["json"] = models.QueryDataWithString(query_string)
}
