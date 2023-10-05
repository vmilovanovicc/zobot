package main

import (
	"time"
	"zobot/backend"
)

func main() {
	// remove these
	backend.TranslateText("hello world", "es")
	//locationURL, _ := backend.CreateBucket("")
	//backend.GetBucketName(locationURL)
	backend.GetTargetVoice("fr")
	taskId, _, _ := backend.GetSpeechSynthesisTaskId("Testing speech to voice functionality", "", "es", "")
	backend.GetSpeechSynthesisTaskStatus(taskId)
	backend.GetPresignedURL("amazon-polly-audio-10012023182532", "bf923301-3f75-4a0a-adcb-ff0d3325xx09.mp3", time.Now().Add(3600*time.Second))
}
