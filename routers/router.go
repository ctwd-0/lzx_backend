package routers

import (
	"lzx_backend/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/table/init", &controllers.TableController{})
    beego.Router("/image/get_image", &controllers.ImageController{})
}
