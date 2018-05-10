package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//向数据库添加新的文件夹
func insertFolder(model_id string, folders []string) string {
	db := S.DB("database")
	data := []bson.M{}
	for _, val := range folders {
		data = append(data, bson.M{"val":val,"key":bson.NewObjectId()})
	}

	err := db.C("folder").Insert(bson.M{
		"model_id": model_id,
		"folders": data,
		"author": "system_init",
		"modified": time.Now(),
		"old":[]bson.M{},
	})

	if err == nil {
		return ""
	} else {
		return "新建失败"
	}
}

//获取指定构件对应的文件夹信息。
//两维数组分别为文件夹名对应的随机id，和文件夹名
func GetFolders(model_id string) ([][]string, string) {
	db := S.DB("database")

	var data bson.M
	err := db.C("folder").Find(bson.M{"model_id":model_id}).One(&data)
	if err == nil {
		folders := [][]string{[]string{},[]string{}}
		for _, val := range data["folders"].([]interface{}) {
			folders[0] = append(folders[0],val.(bson.M)["val"].(string))
			folders[1] = append(folders[1],val.(bson.M)["key"].(bson.ObjectId).Hex())
		}
		return folders, ""
	} else if err.Error() == "not found" {
		reason := insertFolder(model_id, []string{"图纸", "照片", "正射影像"})
		if reason == "" {
			return GetFolders(model_id)
		} else {
			return [][]string{}, reason
		}
	} else {
		return [][]string{}, "数据库错误"
	}
}

//将指定构件的指定文件夹名转换为对应的id
func ConvertName(model_id, folder_name string) (string, string) {
	folders, reason := GetFolders(model_id)

	if reason == "" {
		for idx, val := range folders[0] {
			if val == folder_name {
				return folders[1][idx], ""
			}
		}
		return "", "不存在的类别"
	}
	return "", reason
}
