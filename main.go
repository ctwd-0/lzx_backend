package main

import (
	_ "lzx_backend/routers"
	"lzx_backend/models"
	"github.com/astaxie/beego"
)

func main() {
	defer models.S.Close()
	models.Start_server()
	beego.Run()
}
