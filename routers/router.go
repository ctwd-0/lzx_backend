package routers

import (
	"lzx_backend/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/table/init", &controllers.TableController{}, "get:InitTable")
    beego.Router("/table/update", &controllers.TableController{}, "get:UpdateValue")
    beego.Router("/table/add_column", &controllers.TableController{}, "get:AddColumn")
    beego.Router("/table/remove_column", &controllers.TableController{}, "get:RemoveColumn")
    beego.Router("/table/rename_column", &controllers.TableController{}, "get:RenameColumn")
    //beego.Router("/image/get_image", &controllers.ImageController{})
    beego.Router("/search", &controllers.SearchController{})
    beego.Router("/query/init", &controllers.QueryController{}, "get:InitQuery")
    beego.Router("/query/add", &controllers.QueryController{}, "post:AddQuery")
    beego.Router("/query/get", &controllers.QueryController{}, "get:GetQuery")
    beego.Router("/query/delete", &controllers.QueryController{}, "post:DeleteQuery")
    beego.Router("/filter/init", &controllers.FilterController{}, "get:InitFilter")
    beego.Router("/filter/update", &controllers.FilterController{}, "get:UpdateFilter")
    beego.Router("/filter/delete", &controllers.FilterController{}, "get:DeleteFilter")
    beego.Router("/filter/add", &controllers.FilterController{}, "get:AddFilter")
    beego.Router("/user/login", &controllers.UserController{}, "get:Login")
    beego.Router("/user/logout", &controllers.UserController{}, "get:Logout")
    beego.Router("/admin/add_user", &controllers.AdminController{}, "get:AddUser")
    beego.Router("/admin/init", &controllers.AdminController{}, "get:InitUser")
    beego.Router("/admin/remove_user", &controllers.AdminController{}, "get:DeleteUser")
    beego.Router("/admin/password", &controllers.AdminController{}, "get:ChangePassword")
    beego.Router("/admin/update", &controllers.AdminController{}, "get:UpdateUser")
    beego.Router("/file/upload", &controllers.FileController{}, "post:Upload")
    beego.Router("/file/upload", &controllers.FileController{}, "options:Options")
    beego.Router("/file/get_files", &controllers.FileController{}, "get:GetAll")
    beego.Router("/file/is_ready", &controllers.FileController{}, "get:Ready")
    beego.Router("/file/update_description", &controllers.FileController{}, "post:Update")
}
