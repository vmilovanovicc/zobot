package backend

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	types2 "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"net/url"
	"time"

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

// CreateBucket creates an S3 bucket to store the output of Amazon Polly Synthesis Tasks i.e. audio streams.
func CreateBucket(bucketName string) (location string, err error) {
	cfg := LoadAWSConfig()
	client := s3.NewFromConfig(cfg)

	bucketName = fmt.Sprintf("%s-%s", "amazon-polly-audio", time.Now().Format("01022006150405"))
	bucketLocationConstraint := types2.BucketLocationConstraint(region)

	params := &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types2.CreateBucketConfiguration{
			LocationConstraint: bucketLocationConstraint,
		}, // basic configuration
	}
	ctx := context.TODO()
	resp, err := client.CreateBucket(ctx, params)
	if err != nil {
		log.Fatalf("failed to create s3 bucket, error:  %v\n", err)
		return "", err
	}
	location = *resp.Location
	fmt.Println(location)
	return location, nil
}

// GetBucketName retrieves the S3 bucket name.
func GetBucketName(locationURL string) (outputS3BucketName string, err error) {
	locationURL, _ = CreateBucket("")
	parsedURL, err := url.Parse(locationURL)
	if err != nil {
		log.Fatalf("failed to parse S3 bucket locationURL, error: %v\n", err)
		return "", err
	}
	hostParts := strings.Split(parsedURL.Host, ".")
	if len(hostParts) < 3 {
		log.Fatalf("invalid s3 location URL: %s", locationURL)
	}
	outputS3BucketName = hostParts[0]
	fmt.Println(outputS3BucketName)
	return outputS3BucketName, nil
}

// GetSpeechSynthesisTaskId function creates a synthesis task which converts provided text into an audio stream.
func GetSpeechSynthesisTaskId(text, bucketName, languageCode, targetVoice string) (taskId, objectName string, err error) {
	cfg := LoadAWSConfig()
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
