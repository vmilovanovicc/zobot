package main

import (
	"zobot/backend"
)

func main() {
	backend.TranslateText("hello world", "en")
	backend.GetTargetVoice("es")
}
