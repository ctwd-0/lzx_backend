package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/models"
	"lzx_backend/utils"
	"time"
)

type FolderController struct {
	beego.Controller
}

func returnFolders(model_id, reason string) bson.M {
	var folders [][]string
	if reason == "" {
		folders, reason = models.GetFolders(model_id)
	}
	var files []bson.M
	if reason == "" {
		files, reason = models.AllFiles(model_id, folders[0][0])
	}

	if reason == "" {
		return bson.M{"success":true, "folders": folders[0], "files": files}
	} else {
		return bson.M{"success":false, "reason": reason}
	}
}

func (c *FolderController) GetFolders() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	model_id := c.GetString("model_id")
	if model_id == "" {
		reason = "参数错误"
	}

	c.Data["json"] = returnFolders(model_id, reason)
}

func (c *FolderController) RenameFolder() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	model_id := c.GetString("model_id")
	new_folder_name := c.GetString("new")
	old_folder_name := c.GetString("old")
	author := "system_test"
	if model_id == "" || new_folder_name == "" || old_folder_name == "" || old_folder_name == new_folder_name {
		reason = "参数错误"
	}
	
	if reason == "" {
		reason = renameFolder(old_folder_name, new_folder_name, model_id, author)
	}

	c.Data["json"] = returnFolders(model_id, reason)
}

func (c *FolderController) RemoveFolderAndMove() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	db := models.S.DB("database")
	model_id := c.GetString("model_id")
	new_folder_name := c.GetString("new")
	old_folder_name := c.GetString("old")
	author := "system_test"
	if model_id == "" || new_folder_name == "" || old_folder_name == "" || old_folder_name == new_folder_name {
		reason = "参数错误"
	}
	
	var old_id, new_id string
	if reason == "" {
		old_id, new_id, reason = removeFolder(old_folder_name, new_folder_name, model_id, author)
	}

	if reason == "" && old_id != "" && new_id != ""{
		_, err := db.C("file").UpdateAll(bson.M{"category": old_id},bson.M{"$set":bson.M{"category":new_id}})
		if err != nil {
			reason = "重命名失败"
		}
	}

	c.Data["json"] = returnFolders(model_id, reason)
}

func (c *FolderController) AddFolder() {
	defer c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	reason := ""
	model_id := c.GetString("model_id")
	folder_name := c.GetString("folder_name")
	author := "system_test"
	if model_id == "" || folder_name == "" {
		reason = "参数错误"
	}

	if reason == "" {
		reason = addFolder(folder_name, model_id, author)
	}

	c.Data["json"] = returnFolders(model_id, reason)
}

func addFolder(folder_name, model_id, author string) string {
	folder, reason := models.GetFolders(model_id)
	if reason == "" {
		for _, val := range folder[0] {
			if val == folder_name {
				reason = "重复的名称"
				break
			}
		}
	}

	if reason == "" {
		folder[0] = append(folder[0], folder_name)
		folder[1] = append(folder[1], bson.NewObjectId().Hex())
		updateFolder(folder, model_id, author)
	}

	return reason
}

func removeFolder(old_folder_name, new_folder_name, model_id, author string) (string, string, string) {
	var old_id, new_id string
	folder, reason := models.GetFolders(model_id)
	if reason == "" {
		if len(folder[0]) == 1 {
			return "", "" , "不能删除最后一个选项"
		}
		new_folder := [][]string{[]string{}, []string{}}
		for idx, val := range folder[0] {
			if val == old_folder_name {
				old_id = folder[1][idx]
				continue
			} else if val == new_folder_name {
				new_id = folder[1][idx]
			}
			new_folder[0] = append(new_folder[0], val)
			new_folder[1] = append(new_folder[1], folder[1][idx])
		}
		if old_id == "" {
			return "", "", "未找到旧名字"
		} else if new_id == "" {
			return old_id, "", renameFolder(old_folder_name, new_folder_name, model_id, author)
		} else {
			reason = updateFolder(new_folder, model_id, author)
			return old_id, new_id, reason
		}
	} else {
		return "", "", reason
	}
}

func renameFolder(old_folder_name, new_folder_name, model_id, author string) string {
	folder, reason := models.GetFolders(model_id)

	if reason == "" {
		index := -1
		exist := false
		if reason == "" {
			for idx, val := range folder[0] {
				if val == old_folder_name {
					index = idx
				} else if val == new_folder_name {
					exist = true
				}
			}
			if index == -1 {
				reason = "未找到旧名字"
			} else if exist {
				reason = "重复的名字"
			} else {
				folder[0][index] = new_folder_name
			}
		}
	}
	
	if reason == "" {
		reason = updateFolder(folder, model_id, author)
	}

	return reason
}

func updateFolder(new_folder [][]string, model_id, author string) string {
	db := models.S.DB("database")

	var r bson.M
	err := db.C("folder").Find(bson.M{"model_id":model_id}).One(&r)
	if err == nil {
		data := utils.StringToData(new_folder)
		err := db.C("folder").UpdateId(r["_id"], bson.M{
			"$push":bson.M{"old":bson.M{"folders":r["folders"],"modified":r["modified"],"author":r["author"]}},
			"$set":bson.M{"folders": data, "modified":time.Now(), "author": author},
		})
		if err == nil {
			return ""
		} else {
			return "更新文件夹失败"
		}
	} else {
		return "查找文件夹失败"
	}
}
