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
    beego.Router("/user/update", &controllers.UserController{}, "get:UpdateUser")
    beego.Router("/user/password", &controllers.UserController{}, "get:ChangePassword")
    beego.Router("/user/login", &controllers.UserController{}, "get:Login")
    beego.Router("/user/logout", &controllers.UserController{}, "get:Logout")
    beego.Router("/admin/add_user", &controllers.AdminController{}, "get:AddUser")
    beego.Router("/admin/init", &controllers.AdminController{}, "get:InitUser")
}
