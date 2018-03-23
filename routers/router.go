package routers

import (
	"lzx_backend/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/table/init", &controllers.TableController{})
    beego.Router("/image/get_image", &controllers.ImageController{})
    beego.Router("/search", &controllers.SearchController{})
    beego.Router("/query/init", &controllers.QueryController{}, "get:InitQuery")
    beego.Router("/query/add", &controllers.QueryController{}, "post:AddQuery")
    beego.Router("/query/get", &controllers.QueryController{}, "get:GetQuery")
    beego.Router("/query/delete", &controllers.QueryController{}, "post:DeleteQuery")
    beego.Router("/filter/init", &controllers.FilterController{}, "get:InitFilter")
    beego.Router("/filter/update", &controllers.FilterController{}, "get:UpdateFilter")
    beego.Router("/filter/delete", &controllers.FilterController{}, "get:DeleteFilter")
    beego.Router("/filter/add", &controllers.FilterController{}, "get:AddFilter")
}
