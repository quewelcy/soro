package soro

import (
	"strings"
)

var audioSuffixes = [...]string{
	".mp3", ".MP3",
}

func isAudio(path string) bool {
	for _, suf := range audioSuffixes {
		if strings.HasSuffix(path, suf) {
			return true
		}
	}
	return false
}

func getAudioContent(path string) string {
	return "<audio controls><source src='/file/?p=" + path + "' type='audio/mpeg'></audio>"
}
