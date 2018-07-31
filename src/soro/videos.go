package soro

import "strings"

var videoSuffixes = [...]string{
	".mp4",
}

func isVideo(path string) bool {
	for _, suf := range videoSuffixes {
		if strings.HasSuffix(path, suf) {
			return true
		}
	}
	return false
}

func getVideoContent(path string) string {
	return "<video  width='320px' controls><source src='/file/?p=" + path + "' type='video/mp4'></video>"
}
