package backend

import (
	context "context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"log"
)

// TranslateText detects the source language and translates the text into a desired language.
func TranslateText(text, language string) (string, error) {
	cfg := LoadAWSConfig()
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
		log.Fatalf("failed to translate, error: %v\n", err)
		return "", err
	}
	fmt.Println(*resp.TranslatedText)
	return *resp.TranslatedText, nil
}
