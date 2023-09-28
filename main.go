package main

import "zobot/backend"

func main() {
	// placeholder
	backend.TranslateText("hello world", "en")
	backend.CreateBucket("")
	backend.GetBucketName("")
	backend.GetTargetVoice("es")
	backend.GetSpeechSynthesisTaskId("Testing speech to voice functionality", "", "es", "")
}
