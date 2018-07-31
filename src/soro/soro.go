package soro

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var root = ""
var htmlPath = "/"
var thumbStorage = getHomeDir() + string(os.PathSeparator) + ".soroThumbs"

//StartWeb starts web https://localhost:443
func StartWeb() {
	if len(os.Args) < 2 {
		log.Println("No root path provided")
		return
	}

	root = os.Args[1]
	if root == "" {
		log.Println("Root path is empty")
		return
	}
	log.Println("Root path is", root)

	//todo differ these functions

	if len(os.Args) > 2 {
		thumb := os.Args[2]
		if thumb == "thumb" {
			makeThumbs()
		}
	} else {
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../res/public"))))
		http.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.FormValue("p"))
		})
		http.HandleFunc("/thumb/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, getFullPath(thumbStorage, r.FormValue("p")))
		})
		http.HandleFunc("/", rootHandler)
		err := http.ListenAndServeTLS(":443", "soro-cert.pem", "soro-key.pem", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	dirPath := r.FormValue("p")
	contentPath := r.FormValue("c")
	tmplTitle, _ := template.ParseFiles("../res/template/title.tm")
	tmplfileRow, _ := template.ParseFiles("../res/template/fileRow.tm")
	tmplfileDownload, _ := template.ParseFiles("../res/template/fileDownload.tm")

	contentTitle := ""
	contentClass := "preview-outer-"
	if contentPath == "" {
		contentClass += "invisible"
	} else {
		contentClass += "visible"
		contentTitle = strings.TrimPrefix(
			strings.TrimPrefix(contentPath, dirPath),
			string(os.PathSeparator))
	}

	data := map[string]interface{}{
		"dirPath":      getFileRowLink(dirPath),
		"dir":          template.HTML(readDir(dirPath, contentPath, tmplfileRow, tmplfileDownload)),
		"thumbs":       template.HTML(readDirThumbs(dirPath)),
		"content":      template.HTML(getContent(contentPath)),
		"contentClass": contentClass,
		"contentTitle": contentTitle,
		"contentBack":  getFileRowLink(dirPath),
	}
	tmplTitle.Execute(w, data)
}

func readDir(path, contentPath string, tmplFileRow, tmplfileDownload *template.Template) string {
	if path == "" {
		path = root
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer // todo reuse single bytes buffer

	if path != root {
		data := map[string]interface{}{
			"class":    "",
			"href":     htmlPath + "?p=" + getUpperDir(path),
			"fileName": "/..",
			"download": "",
		}
		tmplFileRow.Execute(&b, data)
	}

	for _, file := range files {
		if skipInList(file.Name()) {
			continue
		}
		var downloadLink template.HTML
		var href string
		var class string

		filePath := getFullPath(path, file.Name())
		addDownload := false
		fileName := file.Name()
		if file.IsDir() {
			href = htmlPath + "?p=" + filePath
			fileName = "/" + fileName
		} else {
			href = htmlPath + "?p=" + path + "&c=" + filePath
			addDownload = true
			if filePath == contentPath {
				class = "current-preview" //fixme move to const
			}
		}

		if addDownload {
			downloadLink = getDownloadLink(filePath, tmplfileDownload)
		}
		data := map[string]interface{}{
			"class":    class,
			"href":     href,
			"fileName": fileName,
			"download": downloadLink,
		}
		tmplFileRow.Execute(&b, data)
	}
	return b.String()
}

func getDownloadLink(filePath string, tmpl *template.Template) template.HTML {
	var b bytes.Buffer
	data := map[string]interface{}{
		"filePath": filePath,
	}
	tmpl.Execute(&b, data)
	return template.HTML(b.String())
	// todo sequential writing into w
}

func skipInList(fname string) bool {
	return isImage(fname) || isVideo(fname)
}

func readDirThumbs(path string) string {
	if path == "" {
		path = root
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	var b bytes.Buffer
	for _, file := range files {
		fileName := file.Name()
		if !file.IsDir() && (isImage(fileName) || isVideo(fileName)) {
			imgPath := getFullPath(path, fileName)
			class := ""
			if isVideo(fileName) {
				fileName += ".png"
				class = "video-mark"
			}
			getThumbTag(&b, fileName, path, imgPath, class)
		}
	}
	return b.String()
}

func getUpperDir(dir string) string {
	separatorIndex := strings.LastIndex(dir, string(os.PathSeparator))
	if separatorIndex > 0 {
		return dir[:separatorIndex]
	}
	return dir
}

func getFullPath(dir, file string) string {
	if strings.HasSuffix(dir, string(os.PathSeparator)) {
		return dir + file
	}
	return dir + string(os.PathSeparator) + file
}

func getContent(path string) string {
	if isImage(path) {
		return getImageContent(path)
	} else if isAudio(path) {
		return getAudioContent(path)
	} else if isVideo(path) {
		return getVideoContent(path)
	} else if isText(path) {
		return getTextContent(path)
	}
	return ""
}

func getFileRowLink(dirPath string) string {
	if dirPath == "" {
		return "/"
	}
	return "/?p=" + dirPath
}
