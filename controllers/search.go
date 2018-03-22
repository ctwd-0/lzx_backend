package controllers

import (
	//"fmt"
	"lzx_backend/models"
	"github.com/astaxie/beego"
)

type SearchController struct {
	beego.Controller
}

func (c *SearchController) Get() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")

	query_string := c.GetString("query")
	//fmt.Println(query_string)

	c.Data["json"] = models.QueryDataWithString(query_string)
}

func (c *SearchController) Post() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	
	query_string := c.GetString("query")
	//fmt.Println(query_string)

	c.Data["json"] = models.QueryDataWithString(query_string)
}
