package routers

import (
	"lzx_backend/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//返回初始化数据表
	beego.Router("/table/init", &controllers.TableController{}, "get:InitTable")
	//更新数据表中的值
	beego.Router("/table/update", &controllers.TableController{}, "get:UpdateValue")
	//增加数据表列
	beego.Router("/table/add_column", &controllers.TableController{}, "get:AddColumn")
	//删除数据表列
	beego.Router("/table/remove_column", &controllers.TableController{}, "get:RemoveColumn")
	//重命名数据表列
	beego.Router("/table/rename_column", &controllers.TableController{}, "get:RenameColumn")
	//按检索条件返回符合条件的数据表行
	beego.Router("/search", &controllers.SearchController{}, "get:Search")
	//按检索条件返回符合条件的数据表行
	beego.Router("/search", &controllers.SearchController{}, "post:Search")
	//返回所有保存的检索条件
	beego.Router("/query/init", &controllers.QueryController{}, "get:InitQuery")
	//增加检索条件
	beego.Router("/query/add", &controllers.QueryController{}, "post:AddQuery")
	//获取指定的检索条件的内容
	beego.Router("/query/get", &controllers.QueryController{}, "get:GetQuery")
	//删除检索条件
	beego.Router("/query/delete", &controllers.QueryController{}, "post:DeleteQuery")
	//获取所有保存的表头选择器
	beego.Router("/filter/init", &controllers.FilterController{}, "get:InitFilter")
	//更新表头选择器
	beego.Router("/filter/update", &controllers.FilterController{}, "get:UpdateFilter")
	//删除表头选择器
	beego.Router("/filter/delete", &controllers.FilterController{}, "get:DeleteFilter")
	//增加表头选择其
	beego.Router("/filter/add", &controllers.FilterController{}, "get:AddFilter")
	//用户登录
	beego.Router("/user/login", &controllers.UserController{}, "get:Login")
	//用户登出
	beego.Router("/user/logout", &controllers.UserController{}, "get:Logout")
	//管理员增加用户
	beego.Router("/admin/add_user", &controllers.AdminController{}, "get:AddUser")
	//管理员获取用户列表
	beego.Router("/admin/init", &controllers.AdminController{}, "get:InitUser")
	//管理员删除用户
	beego.Router("/admin/remove_user", &controllers.AdminController{}, "get:DeleteUser")
	//管理员设置用户密码
	beego.Router("/admin/password", &controllers.AdminController{}, "get:ChangePassword")
	//管理员更新用户状态
	beego.Router("/admin/update", &controllers.AdminController{}, "get:UpdateUser")
	//管理员页面
	beego.Router("/admin", &controllers.AdminController{}, "get:Admin")
	beego.Router("/admin.html", &controllers.AdminController{}, "get:Admin")
	//上传文件
	beego.Router("/file/upload", &controllers.FileController{}, "post:Upload")
	beego.Router("/file/upload", &controllers.FileController{}, "options:Options")
	//获取某个构件某个文件夹下的所有文件
	beego.Router("/file/get_files", &controllers.FileController{}, "get:GetAll")
	//测试上传文件是否完场
	beego.Router("/file/is_ready", &controllers.FileController{}, "get:Ready")
	//修改文件的描述信息
	beego.Router("/file/update_description", &controllers.FileController{}, "post:Update")
	//删除文件
	beego.Router("/file/remove", &controllers.FileController{}, "post:Remove")
	//下载文件
	beego.Router("/file/download", &controllers.FileController{}, "post:Download")
	//获取某个构件对应的文件夹信息
	beego.Router("/folder/init", &controllers.FolderController{}, "get:GetFolders")
	//重命名文件夹
	beego.Router("/folder/rename", &controllers.FolderController{}, "get:RenameFolder")
	//删除文件夹
	beego.Router("/folder/remove", &controllers.FolderController{}, "get:RemoveFolderAndMove")
	//增加文件夹
	beego.Router("/folder/add", &controllers.FolderController{}, "get:AddFolder")
}
