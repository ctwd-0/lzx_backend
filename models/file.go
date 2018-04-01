package models

import (
	"github.com/satori/go.uuid"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"strings"
	"io/ioutil"
	"os"
	"io"
	"fmt"
	"time"
	"crypto/md5"
	"gopkg.in/mgo.v2/bson"
	"lzx_backend/utils"
	"gopkg.in/gographics/imagick.v3/imagick"
)

const preprocess_dest_dir = "E:/20170109/building_viewer/dist/files/"
const upload_dest_dir = "E:/20170109/building_viewer/dist/uploads/"

func InitDbFile(path string) {
	groups, _ := ioutil.ReadDir(path)
	db := S.DB("database")
	c := db.C("file")

	for _, group := range groups {
		categories, _ := ioutil.ReadDir(path + "/" + group.Name())
		for _, cat := range categories {
			files, _ := ioutil.ReadDir(path + "/" + group.Name() + "/" + cat.Name())
			for _, file := range files {
				model_id := group.Name()
				category := cat.Name()
				filename := file.Name()
				path_name := path + "/" + model_id + "/" + category + "/"+ filename
				parts := strings.Split(filename, ".")
				if parts[len(parts) - 2] == "thumbnail" {
					continue
				}
				fmt.Print(path_name)
				m, _ := utils.Thumbnail(path_name)
				if m["success"].(bool) {
					ext := m["ext"].(string)
					ori_md5 := utils.FileMD5(path_name)
					thu_md5 := utils.FileMD5(m["thumbnail_path"].(string))
					ori_saved_as := ori_md5 + "." + ext
					thu_saved_as := thu_md5 + "." + ext
					utils.CopyFile(m["ori_path"].(string),
						preprocess_dest_dir + ori_md5 + "." + ext)
					utils.CopyFile(m["thumbnail_path"].(string),
						preprocess_dest_dir + thu_md5 + "." + ext)
					if m["success"].(bool) {
						mm, reason := makeFileDocHelper("/dist/files/", model_id, category, ori_md5, thu_md5, ori_saved_as, thu_saved_as, filename, "image", "", 
							m["ori_x"].(int), m["ori_y"].(int), m["thu_x"].(int), m["thu_y"].(int))
						if reason == "" {
							delete(mm, "uuid")
							c.Insert(mm)
						} else {
							fmt.Println(" error")
							fmt.Println(reason)
						}
					}
				}
				fmt.Println(" finished")
			}
		}
	}
}

func saveFileHelper(ori_content, thu_content []byte, ext, thu_ext string) (ori_md5, thu_md5, ori_saved_as, thu_saved_as, reason string) {
	reason = ""

	ori_md5 = fmt.Sprintf("%x",md5.Sum(ori_content))
	ori_saved_as = ori_md5 + "." + ext
	ori_path := upload_dest_dir + ori_saved_as
	out_file, err := os.Create(ori_path)
	if err == nil {
		_, err = out_file.Write(ori_content)
		out_file.Close()
		if err != nil {
			reason = "原始文件写入失败"
		}
	} else {
		reason = "创建原始文件失败"
	}

	if reason == "" {
		thu_md5 = fmt.Sprintf("%x",md5.Sum(thu_content))
		thu_saved_as = thu_md5 + "." + thu_ext
		thu_path := upload_dest_dir + thu_saved_as
		out_file, err := os.Create(thu_path)
		if err == nil {
			_, err = out_file.Write(thu_content)
			out_file.Close()
			if err != nil {
				reason = "原始文件写入失败"
			}
		} else {
			reason = "创建原始文件失败"
		}
	}

	return
}

func makeFileDocHelper(prefix, model_id, category, ori_md5, thu_md5, ori_saved_as, thu_saved_as, filename, type_name, uuid string, ox, oy, tx, ty int) (bson.M, string) {
	category_id, reason := ConvertName(model_id, category)
	return bson.M{
		"model_id": model_id,
		"category": category_id,
		"original_md5": ori_md5,
		"thumbnail_md5": thu_md5,
		"original_saved_as": ori_saved_as,
		"thumbnail_saved_as": thu_saved_as,
		"original_path": prefix + ori_saved_as,
		"thumbnail_path": prefix + thu_saved_as,
		"original_name": filename,
		"original_width": ox,
		"original_height": oy,
		"thumbnail_width": tx,
		"thumbnail_height": ty,
		"type": type_name,
		"description": "",
		"uuid": uuid,
		"created": time.Now(),
		"deleted": false,
	}, reason
}

func ProcessPdf(file multipart.File, filename string, model_id string, category string, uuid uuid.UUID) {
	name_parts := strings.Split(filename, ".")
	ext := name_parts[len(name_parts) - 1]

	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	reason := ""
	var err error
	var thu_temp *os.File

	var ori_content, thu_content []byte
	var ori_md5, thu_md5, ori_saved_as, thu_saved_as string
	ori_content, err = ioutil.ReadAll(file)
	if err != nil {
		reason = "读取上传文件失败"
	}

	if reason == "" {
		err = mw.ReadImageBlob(ori_content)
		if err != nil {
			reason = "文件解码失败"
		}
	}

	var thu_temp_name string
	if reason == "" {
		thu_temp, err = ioutil.TempFile("", "mw_temp")
		if err != nil {
			reason = "创建临时文件失败"
		} else {
			thu_temp_name = thu_temp.Name() + ".png"
			thu_temp.Close()
			os.Remove(thu_temp.Name())
			os.Remove(thu_temp_name)
		}
	}
	
	if mw.GetNumberImages() == 0 {
		reason = "文档为空"
	}

	if reason == "" {
		mw.SetIteratorIndex(0)
		err = mw.WriteImage(thu_temp.Name() + ".png")
		if err != nil {
			reason = "写入缩略图失败"
		}
	}

	if reason == "" {
		thu_content, err = ioutil.ReadFile(thu_temp.Name() + ".png")
		if err != nil {
			reason = "读取缩略图失败"
		}
	}

	if reason == "" {
		thu_temp, err = os.Open(thu_temp_name)
		if err != nil {
			reason = "打开临时缩略图文件失败"
		}
	}

 	if reason == "" {
		ori_md5, thu_md5, ori_saved_as, thu_saved_as, reason = saveFileHelper(ori_content, thu_content, ext, "png")
	}

	var thu image.Image
	if reason == "" {
		thu_temp.Seek(0,io.SeekStart)
		thu, err = png.Decode(thu_temp)
		if err != nil {
			reason = "缩略图解码失败"
		}
	}

	var m bson.M
	if reason == "" {
		tm := thu.Bounds().Max
		m, reason = makeFileDocHelper("/dist/uploads/", model_id, category, ori_md5, thu_md5, ori_saved_as, thu_saved_as, filename, "pdf", uuid.String(), tm.X, tm.Y, tm.X, tm.Y)
	}
	if reason == "" {
		c := S.DB("database").C("file")
		delete(m, "original_width")
		delete(m, "original_height")
		c.Insert(m)
	}
}

func ProcessUploadedFile(file multipart.File, filename string, model_id string, category string, uuid uuid.UUID) {
	name_parts := strings.Split(filename, ".")
	ext := name_parts[len(name_parts) - 1]
	l_ext := strings.ToLower(ext)
	
	if l_ext == "pdf" {
		ProcessPdf(file, filename, model_id, category, uuid)
		return
	}

	ori, thu, reason := utils.ThumbnailFile(file, l_ext)

	var thu_temp *os.File
	var err error
	if reason == "" {
		thu_temp, err = ioutil.TempFile("", "image_")
		if err != nil {
			reason = "创建临时文件失败"
		}
	}

	if reason == "" {
		if l_ext == "jpg" || l_ext == "jpeg" {
			err = jpeg.Encode(thu_temp, thu, &jpeg.Options{100})
		} else if l_ext == "png" {
			err = png.Encode(thu_temp, thu)
		}
		if err != nil {
			reason = err.Error()
		}
	}

	var ori_content, thu_content []byte
	var ori_md5, thu_md5, ori_saved_as, thu_saved_as string
	if reason == "" {
		file.Seek(0,io.SeekStart)
		ori_content, err = ioutil.ReadAll(file)
		if err != nil {
			reason = "读取上传文件失败"
		}
	}

	if reason == "" {
		thu_temp.Seek(0,io.SeekStart)
		thu_content, err = ioutil.ReadAll(thu_temp)
		if err != nil {
			reason = "读取缩略图文件失败"
		}
		thu_temp.Close()
		os.Remove(thu_temp.Name())
	}

 	if reason == "" {
		ori_md5, thu_md5, ori_saved_as, thu_saved_as, reason = saveFileHelper(ori_content, thu_content, ext, "png")
	}

	var m bson.M
	if reason == "" {

		tm := thu.Bounds().Max
		om := ori.Bounds().Max
		m, reason = makeFileDocHelper("/dist/uploads/", model_id, category, ori_md5, thu_md5, ori_saved_as, thu_saved_as, filename, "image", uuid.String(), om.X, om.Y, tm.X, tm.Y)

	}
	if reason == "" {
		c :=  S.DB("database").C("file")
		c.Insert(m)
	}
}

func AllFiles(model_id, category string) ([]bson.M, string){
	db := S.DB("database")
	category_id, reason := ConvertName(model_id, category)

	data:= []bson.M{}
	if reason == "" {
		err := db.C("file").Find(bson.M{
			"model_id": model_id, "category": category_id, "deleted": false,
		}).Select(bson.M{
			"deleted":0,"created":0,"uuid":0,
			"original_md5":0,"thumbnail_md5":0,"thumbnail_saved_as":0,"original_saved_as":0,
		}).Sort("-created").All(&data)
		if err != nil {
			reason = "数据库错误"
		}
	}
	return data, reason
}
