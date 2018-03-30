package main

import (
	_ "lzx_backend/routers"
	"lzx_backend/models"
	"github.com/astaxie/beego"
)

func main() {
	defer models.S.Close()
	models.Start_server()
	//models.InitDbFile("E:/images")
	//models.InitDbTable("E:/20170109/数据表.csv")
	beego.Run()
}
