package backend

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
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
func GetSpeechSynthesisTaskId(text, bucketName, languageCode, targetVoice string) (string, string, error) {
	cfg := LoadAWSConfig()
	client := polly.NewFromConfig(cfg)

	targetVoice, _ = GetTargetVoice(languageCode)
	targetVoiceId := types.VoiceId(targetVoice)
	bucketLocation, _ := CreateBucket("")
	bucketName, _ = GetBucketName(bucketLocation)

	params := &polly.StartSpeechSynthesisTaskInput{
		Text:               aws.String(text),
		OutputS3BucketName: aws.String(bucketName),
		OutputFormat:       types.OutputFormatMp3,
		Engine:             types.EngineStandard,
		VoiceId:            targetVoiceId,
	}

	ctx := context.TODO()
	resp, err := client.StartSpeechSynthesisTask(ctx, params)
	if err != nil {
		log.Fatalf("failed to retrieve task id and object name, error:  %v\n", err)
		return "", "", err
	}
	uriParts := strings.Split(*resp.SynthesisTask.OutputUri, "/")
	objectName := uriParts[len(uriParts)-1]
	fmt.Println(*resp.SynthesisTask.TaskId)
	fmt.Println(objectName)
	return *resp.SynthesisTask.TaskId, objectName, nil
}
