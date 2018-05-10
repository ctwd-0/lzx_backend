package main

import (
	_ "lzx_backend/routers"
	"lzx_backend/models"
	"github.com/astaxie/beego"
	"net/http"
	"github.com/astaxie/beego/context"
)

// 服务器入口
func main() {
	defer models.S.Close()

	// 建立数据库链接
	models.Start_server()

	//初始化图片文件并建立数据库内容
	//models.InitDbFile("E:/images")
	
	//初始化数据表并建立数据库内容
	//models.InitDbTable("E:/20170109/数据表.csv")

	//设置静态文件夹，并对部分请求设置转发
	beego.BConfig.WebConfig.StaticDir["/dist"] = "dist"
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)

	//启动服务器
	beego.Run()
}

//对部分请求进行转发
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
