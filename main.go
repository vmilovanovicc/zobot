package main

import (
	"zobot/backend"
)

func main() {
	// placeholder
	backend.TranslateText("hello world", "en")
	backend.GetTargetVoice("es")
	backend.GetSpeechSynthesisTaskId("Testing speech to voice functionality", "bucket-polly19", "es", "")
}
