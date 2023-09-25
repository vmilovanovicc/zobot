package backend

import (
	context "context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"log"
)

// TranslateText detects the language and translates the text into a desired language.
func TranslateText(text, language string) (string, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatalf("failed to load configuration, %v\n", err)
	}

	client := translate.NewFromConfig(cfg)

	// If auto is specified, Amazon Translate will call Amazon Comprehend to detect the source language.
	sourceLanguageCode := "auto"

	params := &translate.TranslateTextInput{
		Text:               &text,
		SourceLanguageCode: &sourceLanguageCode,
		TargetLanguageCode: &language,
	}
	ctx := context.TODO()
	resp, err := client.TranslateText(ctx, params)
	if err != nil {
		log.Fatalf("cannot translate, %v\n", err)
		return "", err
	}
	//println(*resp.TranslatedText)
	return *resp.TranslatedText, nil
}
