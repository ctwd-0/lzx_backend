package utils

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"image"
	"os"
	"strings"
)

func Thumbnail(path string) (map[string]interface{}, string) {
	reason := ""
	m := map[string]interface{}{}
	parts := strings.Split(path, "/")
	name := parts[len(parts) - 1]
	name_parts := strings.Split(name, ".")
	ext := name_parts[len(name_parts) - 1]

	file, err := os.Open(path)
	if err != nil {
		reason = "打开文件失败"
	}

	var ori_image image.Image
	err = nil
	if reason == "" {
		if ext == "jpg" || ext == "jpeg"{
			ori_image, err = jpeg.Decode(file)
		} else if ext == "png" {
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
	if reason == "" {
		thumbnail := resize.Thumbnail(1280, 720, ori_image, resize.Lanczos3)
		pos := strings.LastIndex(path, ".")
		thumbnail_path = path[0:pos] + ".thumbnail." + ext
		//thumbnail_path = "thumbnail." + ext
		out, err := os.Create(thumbnail_path)
		if err != nil {
			reason = "创建文件失败"
		} else {
			if ext == "jpg" || ext == "jpeg" {
				err = jpeg.Encode(out, thumbnail, &jpeg.Options{100})
			} else if ext == "png" {
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
		m["ori_path"] = path
		m["thumbnail_path"] = thumbnail_path
		m["ext"] = ext
	} else {
		m["success"] = false
		m["reason"] = reason
	}
	return m, reason
}
