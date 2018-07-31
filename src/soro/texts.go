package soro

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/russross/blackfriday"
)

var dotMd = ".md"
var textSuffixes = [...]string{
	".txt", ".ini", dotMd,
}

func isText(path string) bool {
	for _, suf := range textSuffixes {
		if strings.HasSuffix(path, suf) {
			return true
		}
	}
	return false
}

func getTextContent(path string) string {
	if strings.HasSuffix(path, dotMd) {
		return readMd(path)
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	return string(bytes)
}

func readMd(path string) string {
	lind := strings.LastIndex(path, string(os.PathSeparator))
	npath := "/file/?p=" + strings.Replace(path[0:lind+1], `\`, `\\`, -1)
	f, _ := ioutil.ReadFile(path)
	sf := strings.Replace(string(f), "](", "]("+npath, -1)
	md := blackfriday.MarkdownBasic([]byte(sf))
	return string(md)
}
