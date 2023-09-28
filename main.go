package main

import (
	"zobot/backend"
)

func main() {
	// placeholder
	backend.TranslateText("hello world", "en")
	backend.CreateBucket("amazon-polly-output-bkt")
	backend.GetTargetVoice("es")
	backend.GetSpeechSynthesisTaskId("Testing speech to voice functionality", "amazon-polly-output-bkt", "es", "")
}
