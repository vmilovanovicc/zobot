package backend

import (
	"fmt"
	"log"
)

// GetTargetVoice connects the Amazon Translate language code with the Amazon Polly voice.
func GetTargetVoice(language string) (string, error) {
	langToPollyVoice := map[string]string{
		"en": "Emma",
		"es": "Lucia",
		"fr": "Mathieu",
	}
	targetVoice, found := langToPollyVoice[language]
	if !found {
		err := fmt.Errorf("voice mapping not found for the language: %s\n", language)
		log.Printf("%v", err)
		return "", err
	}
	fmt.Println(targetVoice)
	return targetVoice, nil
}
