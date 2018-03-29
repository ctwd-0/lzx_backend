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
)

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

	var ori_image image.Image
	err = nil
	if reason == "" {
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
		file.Close()
	}

	var thumbnail_path string
	var thumbnail image.Image
	if reason == "" {
		thumbnail = resize.Thumbnail(1280, 720, ori_image, resize.Lanczos3)
		pos := strings.LastIndex(path, ".")
		thumbnail_path = path[0:pos] + ".thumbnail." + ext
		//thumbnail_path = "thumbnail." + ext
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
			defer out.Close()
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
		m["thumbnail_x"] = thumbnail.Bounds().Max.X
		m["thumbnail_y"] = thumbnail.Bounds().Max.Y
		//m["thumbnail"] = thumbnail
		m["ext"] = ext
	} else {
		m["success"] = false
		m["reason"] = reason
	}
	return m, reason
}

func FileMD5(path string) string {
	file, _ := os.Open(path)
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	return fmt.Sprintf("%x",md5.Sum(body))
}

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
