package models

import (
	"github.com/satori/go.uuid"
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
				file_name := file.Name()
				path_name := path + "/" + model_id + "/" + category + "/"+ file_name
				parts := strings.Split(file_name, ".")
				if parts[len(parts) - 2] == "thumbnail" {
					continue
				}
				fmt.Print(path_name)
				m, _ := utils.Thumbnail(path_name)
				if m["success"].(bool) {
					ext := m["ext"].(string)
					ori_md5 := utils.FileMD5(path_name)
					thumbnail_md5 := utils.FileMD5(m["thumbnail_path"].(string))
					utils.CopyFile(m["ori_path"].(string),
						preprocess_dest_dir + ori_md5 + "." + ext)
					utils.CopyFile(m["thumbnail_path"].(string),
						preprocess_dest_dir + thumbnail_md5 + "." + ext)
					if m["success"].(bool) {
						c.Insert(bson.M{
							"model_id": model_id,
							"category": category,
							"original_md5": ori_md5,
							"thumbnail_md5": thumbnail_md5,
							"original_saved_as": ori_md5 + "." + ext,
							"thumbnail_saved_as": thumbnail_md5 + "." + ext,
							"original_path": "/dist/files/" + ori_md5 + "." + ext,
							"thumbnail_path": "/dist/files/" + thumbnail_md5 + "." + ext,
							"original_name": file_name,
							"thumbnail_width": m["thumbnail_x"],
							"thumbnail_height": m["thumbnail_y"],
							"original_width": m["thumbnail_x"],
							"original_height": m["thumbnail_y"],
							"type": "image",
							"created": time.Now(),
							"deleted": false,
						})
					}
				}
				fmt.Println(" finished")
			}
		}
	}
}

func ProcessUploadedFile(file multipart.File, filename string, model_id string, category string, uuid uuid.UUID) {
	name_parts := strings.Split(filename, ".")
	ext := name_parts[len(name_parts) - 1]
	l_ext := strings.ToLower(ext)

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
	var ori_md5, thu_md5, ori_saved_as, thu_saved_as, ori_path, thu_path string
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
		ori_md5 = fmt.Sprintf("%x",md5.Sum(ori_content))
		ori_saved_as = ori_md5 + "." + ext
		ori_path = upload_dest_dir + ori_saved_as
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
	}

	if reason == "" {
		thu_md5 = fmt.Sprintf("%x",md5.Sum(thu_content))
		thu_saved_as = thu_md5 + "." + ext
		thu_path = upload_dest_dir + thu_saved_as
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

	if reason == "" {
		db := S.DB("database")
		c := db.C("file")
		c.Insert(bson.M{
			"model_id": model_id,
			"category": category,
			"original_md5": ori_md5,
			"thumbnail_md5": thu_md5,
			"original_saved_as": ori_saved_as,
			"thumbnail_saved_as": thu_saved_as,
			"original_path": "/dist/uploads/" + ori_saved_as,
			"thumbnail_path": "/dist/uploads/" + thu_saved_as,
			"original_name": filename,
			"thumbnail_width": thu.Bounds().Max.X,
			"thumbnail_height": thu.Bounds().Max.Y,
			"original_width": ori.Bounds().Max.X,
			"original_height": ori.Bounds().Max.Y,
			"type": "image",
			"uuid": uuid.String(),
			"created": time.Now(),
			"deleted": false,
		})
	}
}
