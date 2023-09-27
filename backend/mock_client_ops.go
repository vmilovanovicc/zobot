package backend

// For unit testing, mock out the SDK.
// Using Go interfaces, instead of concrete service client.

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/translate"
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
		Text:               &text,
		TargetLanguageCode: &language,
		SourceLanguageCode: &sourceLanguageCode,
	})
	if err != nil {
		return "", err
	}
	return *response.TranslatedText, nil
}
