package main

import (
	_ "lzx_backend/routers"
	"lzx_backend/models"
	"github.com/astaxie/beego"
	"net/http"
	"github.com/astaxie/beego/context"
)

func main() {
	defer models.S.Close()
	models.Start_server()
	//models.InitDbFile("E:/images")
	//models.InitDbTable("E:/20170109/数据表.csv")
	beego.BConfig.WebConfig.StaticDir["/dist"] = "dist"
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
	beego.Run()
}

func TransparentStatic(ctx *context.Context) {
	path := ctx.Request.URL.Path
	static_path := ""
	if path == "/" {
		static_path = "dist/index.html"
	} else if path == "/index" || path == "index" || path == "/index.html" {
		static_path = "dist/index.html"
	} else if path == "/viewer" || path == "viewer" || path == "/viewer.html" {
		static_path = "dist/viewer.html"
	}
	if static_path != "" {
		http.ServeFile(ctx.ResponseWriter, ctx.Request, static_path)
	}
}
