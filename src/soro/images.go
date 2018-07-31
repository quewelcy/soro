package soro

import (
	"bytes"
	"html/template"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

const (
	previewWidth = 500
)

var pngSuffixes = []string{".png", ".PNG"}
var jpgSuffixes = []string{".jpg", ".JPG", ".jpeg", ".JPEG"}
var imgSuffixes = append(pngSuffixes, jpgSuffixes...)

func isImage(path string) bool {
	for _, suf := range imgSuffixes {
		if strings.HasSuffix(path, suf) {
			return true
		}
	}
	return false
}

func getImageContent(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return ""
	}

	defer file.Close()

	var w, h int
	if hasSuffix(pngSuffixes, path) {
		image, _ := png.DecodeConfig(file)
		w, h = image.Width, image.Height
	} else if hasSuffix(jpgSuffixes, path) {
		image, _ := jpeg.DecodeConfig(file)
		w, h = image.Width, image.Height
	}
	previewHeight := previewWidth * h / w
	data := map[string]interface{}{
		"path":          path,
		"previewWidth":  previewWidth,
		"previewHeight": previewHeight,
	}
	var b bytes.Buffer
	tmpl, _ := template.ParseFiles("../res/template/img.tm")
	tmpl.Execute(&b, data)
	return b.String()
}

func hasSuffix(suffixes []string, path string) bool {
	for _, suf := range suffixes {
		if strings.HasSuffix(path, suf) {
			return true
		}
	}
	return false
}

func getThumbTag(b *bytes.Buffer, thumbName, dirPath, imgPath, class string) {
	tmplThumb, _ := template.ParseFiles("../res/template/thumb.tm")
	data := map[string]interface{}{
		"dirPath":   dirPath,
		"imgPath":   imgPath,
		"thumbName": thumbName,
		"class":     class,
	}
	tmplThumb.Execute(b, data)
}

func haveThumb(thumbPath string) bool {
	_, err := os.Open(thumbPath)
	return err == nil
}
