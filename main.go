package main

import "zobot/backend"

func main() {
	// placeholder
	backend.TranslateText("hello world", "en")
	locationURL, _ := backend.CreateBucket("")
	backend.GetBucketName(locationURL)
	backend.GetTargetVoice("es")
	backend.GetSpeechSynthesisTaskId("Testing speech to voice functionality", "", "es", "")
}
