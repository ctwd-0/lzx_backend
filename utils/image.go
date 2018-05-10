package utils

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"image"
	"os"
	"fmt"
	"crypto/md5"
	"io/ioutil"
	"io"
	"strings"
	"mime/multipart"
)

//将file对应的文件，按照l_ext制定的格式解析为图片，并生成缩略图。
//返回原图，缩略图，以及错误信息。
func ThumbnailFile(file multipart.File, l_ext string) (image.Image, image.Image, string) {
	var ori_image image.Image
	var err error
	reason := ""
	if l_ext == "jpg" || l_ext == "jpeg" {
		ori_image, err = jpeg.Decode(file)
	} else if l_ext == "png" {
		ori_image, err = png.Decode(file)
	} else {
		reason = "不支持的文件格式"
	}
	if err != nil {
		reason = "文件解码错误"
	}

	var thumbnail image.Image
	if reason == "" {
		thumbnail = resize.Thumbnail(1280, 720, ori_image, resize.Lanczos3)
	}

	return ori_image, thumbnail, reason
}

// 将path对应的文件进行缩略图操作。返回包含所有信息的map和错误信息
func Thumbnail(path string) (map[string]interface{}, string) {
	reason := ""
	m := map[string]interface{}{}
	parts := strings.Split(path, "/")
	name := parts[len(parts) - 1]
	name_parts := strings.Split(name, ".")
	ext := name_parts[len(name_parts) - 1]
	l_ext := strings.ToLower(ext)

	file, err := os.Open(path)
	if err != nil {
		reason = "打开文件失败"
	}

	var ori_image, thumbnail image.Image
	if reason == "" {
		ori_image, thumbnail, reason = ThumbnailFile(file, l_ext)
	}
	file.Close()

	var thumbnail_path string
	if reason == "" {
		pos := strings.LastIndex(path, ".")
		thumbnail_path = path[0:pos] + ".thumbnail." + ext
		out, err := os.Create(thumbnail_path)
		if err != nil {
			reason = "创建文件失败"
		} else {
			if l_ext == "jpg" || l_ext == "jpeg" {
				err = jpeg.Encode(out, thumbnail, &jpeg.Options{100})
			} else if l_ext == "png" {
				err = png.Encode(out, thumbnail)
			}
			if err != nil {
				reason = err.Error()
			}
			out.Close()
		}
	}

	if reason == "" {
		m["success"] = true
		//m["ori"] = ori_image
		m["ori_path"] = path
		m["ori_size"] = ori_image.Bounds().Max
		m["ori_x"] = ori_image.Bounds().Max.X
		m["ori_y"] = ori_image.Bounds().Max.Y
		m["thumbnail_size"] = thumbnail.Bounds().Max
		m["thumbnail_path"] = thumbnail_path
		m["thu_x"] = thumbnail.Bounds().Max.X
		m["thu_y"] = thumbnail.Bounds().Max.Y
		//m["thumbnail"] = thumbnail
		m["ext"] = ext
	} else {
		m["success"] = false
		m["reason"] = reason
	}
	return m, reason
}

//返回path对应文件的md5
func FileMD5(path string) string {
	file, _ := os.Open(path)
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	return fmt.Sprintf("%x",md5.Sum(body))
}

//返回file对应文件的md5
func FileMD5File(file multipart.File) string {
	body, _ := ioutil.ReadAll(file)
	return fmt.Sprintf("%x",md5.Sum(body))
}

//将srcName对应的文件拷贝到dstName对应的位置
func CopyFile(srcName, dstName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()
    dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0777)
    if err != nil {
        return
    }
    defer dst.Close()
    return io.Copy(dst, src)
}
