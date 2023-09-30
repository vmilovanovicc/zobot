package main

import "zobot/backend"

func main() {
	// placeholder
	backend.TranslateText("hello world", "es")
	locationURL, _ := backend.CreateBucket("")
	backend.GetBucketName(locationURL)
	backend.GetTargetVoice("fr")
	taskId, _, _ := backend.GetSpeechSynthesisTaskId("Testing speech to voice functionality", "", "es", "")
	backend.GetSpeechSynthesisTaskStatus(taskId)
}
