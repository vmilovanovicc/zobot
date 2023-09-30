package backend

// For unit testing, mock out the SDK.
// Using Go interfaces, instead of concrete service client.

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	typesPolly "github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	typesS3 "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"strings"
)

//////////////////////////////////////////
//         TranslateText Mock
//////////////////////////////////////////

type TranslateTranslateTextAPI interface {
	TranslateText(ctx context.Context, params *translate.TranslateTextInput, optFns ...func(*translate.Options)) (*translate.TranslateTextOutput, error)
}

func TranslateTextFromTranslate(ctx context.Context, api TranslateTranslateTextAPI, text, language string) (string, error) {
	sourceLanguageCode := "auto"
	response, err := api.TranslateText(ctx, &translate.TranslateTextInput{
		Text:               aws.String(text),
		TargetLanguageCode: aws.String(language),
		SourceLanguageCode: aws.String(sourceLanguageCode),
	})
	if err != nil {
		return "", err
	}
	return *response.TranslatedText, nil
}

//////////////////////////////////////////
//         CreateBucket Mock
//////////////////////////////////////////

type S3CreateBucketAPI interface {
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
}

func CreateBucketFromS3(ctx context.Context, api S3CreateBucketAPI, bucketName string) (string, error) {
	response, err := api.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &typesS3.CreateBucketConfiguration{
			LocationConstraint: typesS3.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		return "", err
	}
	return *response.Location, nil
}

//////////////////////////////////////////
//    StartSpeechSynthesisTask Mock
//////////////////////////////////////////

type PollyStartSpeechSynthesisTaskAPI interface {
	StartSpeechSynthesisTask(ctx context.Context, params *polly.StartSpeechSynthesisTaskInput, optFns ...func(*polly.Options)) (*polly.StartSpeechSynthesisTaskOutput, error)
}

func StartSpeechSynthesisTaskFromPolly(ctx context.Context, api PollyStartSpeechSynthesisTaskAPI, text, bucketName, languageCode, targetVoice string) (string, string, error) {
	response, err := api.StartSpeechSynthesisTask(ctx, &polly.StartSpeechSynthesisTaskInput{
		Text:               aws.String(text),
		OutputS3BucketName: aws.String(bucketName),
		VoiceId:            typesPolly.VoiceId(targetVoice),
		LanguageCode:       typesPolly.LanguageCode(languageCode),
		OutputFormat:       typesPolly.OutputFormatMp3,
		Engine:             typesPolly.EngineStandard,
	})
	if err != nil {
		return "", "", err
	}
	outputURI := strings.Split(*response.SynthesisTask.OutputUri, "/")
	objectName := outputURI[len(outputURI)-1]
	return *response.SynthesisTask.TaskId, objectName, nil
}
