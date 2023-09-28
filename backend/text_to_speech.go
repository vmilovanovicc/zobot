package backend

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	"log"
	"strings"
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

// GetSpeechSynthesisTaskId function creates a synthesis task which converts provided text into an audio stream.
func GetSpeechSynthesisTaskId(text, bucketName, languageCode, targetVoice string) (taskId, objectName string, err error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatalf("failed to load configuration, %v\n", err)
	}

	client := polly.NewFromConfig(cfg)

	outputFormat := types.OutputFormatMp3
	engine := types.EngineStandard
	targetVoice, _ = GetTargetVoice(languageCode)
	targetVoiceId := types.VoiceId(targetVoice)

	params := &polly.StartSpeechSynthesisTaskInput{
		Text:               &text,
		OutputS3BucketName: &bucketName,
		OutputFormat:       outputFormat,
		Engine:             engine,
		VoiceId:            targetVoiceId,
	}

	ctx := context.TODO()
	resp, err := client.StartSpeechSynthesisTask(ctx, params)
	if err != nil {
		log.Fatalf("failed to retrieve task id and object name, error:  %v\n", err)
		return "", "", err
	}
	taskId = *resp.SynthesisTask.TaskId
	uriParts := strings.Split(*resp.SynthesisTask.OutputUri, "/")
	objectName = uriParts[len(uriParts)-1]
	fmt.Println(taskId)
	fmt.Println(objectName)
	return taskId, objectName, nil
}
