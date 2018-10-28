package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/nfnt/resize"
)

const (
	thumbHeight = 80
	thumbWidth  = int(thumbHeight * 1.75)
)

func startThumbMaker(rootPath, storagePath string) {
	err := os.MkdirAll(storagePath, 0700) // bug permission denied
	if err != nil {
		log.Println(err)
	}
	makeThumbsOfPath(rootPath, storagePath)
}

func makeThumbsOfPath(rootPath, storagePath string) {
	if rootPath == "" || storagePath == "" {
		return
	}

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fullPath := getFullPath(rootPath, file.Name())
		if file.IsDir() {
			makeThumbsOfPath(fullPath, storagePath)
		} else if isImage(file.Name()) {
			thumbPath := getFullPath(storagePath, file.Name())
			if !isThumbExists(thumbPath) {
				createImgThumb(fullPath, thumbPath)
			}
		} else if isVideo(file.Name()) {
			thumbPath := getFullPath(storagePath, file.Name()+".png")
			if !isThumbExists(thumbPath) {
				createVideoThumb(fullPath, thumbPath)
			}
			//todo duplicates
		}
	}
}

func isThumbExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func createImgThumb(path, thumbPath string) {
	log.Println("Create thumb", path, "->", thumbPath)

	file, err := os.Open(path)
	if err != nil {
		log.Println("Error reading img", path, err)
		return
	}
	defer file.Close()

	var img image.Image
	if hasSuffix(pngSuffixes, path) {
		img, _ = png.Decode(file)
	} else if hasSuffix(jpgSuffixes, path) {
		img, _ = jpeg.Decode(file)
	}

	var thumb image.Image
	if img.Bounds().Max.Y <= thumbHeight {
		thumb = img
	} else {
		thumb = resize.Resize(0, thumbHeight, img, resize.Lanczos3)
	}
	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		log.Println("Error creating thumb", thumbPath, err)
	}
	defer thumbFile.Close()
	err = jpeg.Encode(thumbFile, thumb, nil)
	if err != nil {
		log.Println("Error encoding thumb", thumbPath, err)
	}
}

func createVideoThumb(path, thumbPath string) {
	log.Println("Create video thumb", path, "->", thumbPath)

	out, err := exec.Command("ffmpeg", "-i", path, "-y", "-an", "-ss", "00:00:02", "-vcodec", "png", "-r", "1", "-vframes", "1", "-s", strconv.Itoa(thumbWidth)+"X"+strconv.Itoa(thumbHeight), thumbPath).Output()
	if err != nil {
		log.Println("Error creating video thumb", err)
	} else {
		log.Println(out)
	}
}
