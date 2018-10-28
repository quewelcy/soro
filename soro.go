package main

import (
	"bufio"
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
)

var homeDir string

var rootPath = ""
var htmlPath = "/"
var resPath = ""
var thumbsPath = ""

func main() {
	props := readConf()
	switch props["method"] {
	case "web":
		startWeb(props["port"], props["root"], props["thumbs"], props["resources"], props["cert"], props["key"])
	case "thumbs":
		startThumbMaker(props["root"], props["thumbs"])
	}
}

func readConf() map[string]string {
	if len(os.Args) == 1 {
		log.Fatal("No config provided")
	}
	configPath := os.Args[1]
	if configPath == "" {
		log.Fatal("Config path is empty")
	}
	configPath = resolveFullPath(configPath)
	log.Println("Config path is", configPath)

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal("No config file found")
	}
	props := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		prop := strings.TrimSpace(scanner.Text())
		if prop == "" || prop[0] == '#' {
			continue
		}
		kv := strings.Split(prop, "=")
		props[kv[0]] = resolveFullPath(kv[1])
	}
	return props
}

func resolveFullPath(path string) string {
	if strings.HasPrefix(path, "~") {
		return getHomeDir() + path[1:]
	}
	return path
}

func getHomeDir() string {
	if homeDir == "" {
		usr, _ := user.Current()
		homeDir = usr.HomeDir
	}
	return homeDir
}

func startWeb(port, root, thumbs, res, certPath, keyPath string) {
	if port == "" {
		log.Fatal("Port is not set")
	}
	if root == "" {
		log.Fatal("Root path is not set")
	}
	if thumbs == "" {
		log.Fatal("Thumb storage path is not set")
	}
	rootPath = root
	thumbsPath = thumbs
	resPath = res
	log.Println("Root path is", rootPath)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(resPath+"/public"))))
	http.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.FormValue("p"))
	})
	http.HandleFunc("/thumb/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, getFullPath(thumbsPath, r.FormValue("p")))
	})
	http.HandleFunc("/", rootHandler)

	var err error
	if certPath != "" && keyPath != "" {
		err = http.ListenAndServeTLS(port, certPath, keyPath, nil)
	} else {
		err = http.ListenAndServe(port, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func haveAccess(r *http.Request) bool {
	hash := r.Header.Get("hash")
	if hash == "good" {
		return true
	}
	return false
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	dirPath := r.FormValue("p")
	contentPath := r.FormValue("c")
	tmplTitle, _ := template.ParseFiles(resPath + "/template/title.tm")
	tmplfileRow, _ := template.ParseFiles(resPath + "/template/fileRow.tm")
	tmplfileDownload, _ := template.ParseFiles(resPath + "/template/fileDownload.tm")

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

func readDir(path, contentPath string, tmplFileRow, tmplFileDownload *template.Template) string {
	if path == "" {
		path = rootPath
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer // todo reuse single bytes buffer

	if path != rootPath {
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
			downloadLink = getDownloadLink(filePath, tmplFileDownload)
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
		path = rootPath
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
