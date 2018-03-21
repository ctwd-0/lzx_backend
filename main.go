package main

import (
	_ "lzx_backend/routers"
	"lzx_backend/models"
	"github.com/astaxie/beego"
)

func main() {
	models.Start_server()
	beego.Run()
	defer models.S.Close()
}
