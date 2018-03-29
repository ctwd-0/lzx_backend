package controllers

// import (
// 	"github.com/astaxie/beego"
// 	"gopkg.in/mgo.v2/bson"
// 	"lzx_backend/models"
// )

// type ImageController struct {
// 	beego.Controller
// }

// func (c *ImageController) Get() {
// 	defer c.ServeJSON()
// 	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")   
// 	db := models.S.DB("database")
// 	model_id := c.GetString("model_id")
// 	cat_index, err := c.GetInt("cat_index")
// 	m := make(map[string]interface{})

// 	if(model_id == "" || err != nil) {
// 		m["success"] = false
// 		m["reason"] = "arg"
// 		c.Data["json"] = m
// 		return 
// 	}

// 	var data interface{}
// 	err = db.C("images").Find(bson.M{"model_id":model_id}).Select(bson.M{"_id":0}).One(&data)
// 	if err == nil {
// 		if len(data.(bson.M)["category"].([]interface{})) > cat_index {
// 			cat_name := data.(bson.M)["category"].([]interface{})[cat_index].(string)
// 			m["content"] = data.(bson.M)[cat_name]
// 		} else {
// 			m["success"] = false
// 			m["reason"] = "out of range"
// 		}
// 	} else {
// 		m["success"] = false
// 		m["reason"] = "db"
// 	}
	
// 	c.Data["json"] = m
// }
