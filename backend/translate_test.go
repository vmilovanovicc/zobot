package backend

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"testing"
)

type mockTranslateTextAPI func(ctx context.Context, params *translate.TranslateTextInput, optFns ...func(*translate.Options)) (*translate.TranslateTextOutput, error)

func (m mockTranslateTextAPI) TranslateText(ctx context.Context, params *translate.TranslateTextInput, optFns ...func(*translate.Options)) (*translate.TranslateTextOutput, error) {
	return m(ctx, params, optFns...)
}

func TestTranslateTextFromTranslate(t *testing.T) {
	cases := []struct {
		client               func(t *testing.T) TranslateTranslateTextAPI
		language             string
		text                 string
		expectTranslatedText string
	}{
		{
			client: func(t *testing.T) TranslateTranslateTextAPI {
				return mockTranslateTextAPI(func(ctx context.Context, params *translate.TranslateTextInput, optFns ...func(*translate.Options)) (*translate.TranslateTextOutput, error) {
					t.Helper()
					if params.Text == nil {
						t.Fatal("expect text to not be nil")
					}
					if params.TargetLanguageCode == nil {
						t.Fatal("expect language code not to be nil")
					}
					if len(*params.TargetLanguageCode) > 5 {
						t.Fatal("expect language code to have length less than or equal 5")
					}
					if e, a := "es", *params.TargetLanguageCode; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}

					return &translate.TranslateTextOutput{
						TranslatedText: aws.String("Buenas noches"),
					}, nil
				})
			},
			language:             "es",
			text:                 "Good evening",
			expectTranslatedText: "Buenas noches",
		},
	}

	for _, tt := range cases {
		t.Run("Test translation", func(t *testing.T) {
			ctx := context.TODO()
			gotTranslatedText, err := TranslateTextFromTranslate(ctx, tt.client(t), tt.text, tt.language)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if tt.expectTranslatedText != gotTranslatedText {
				t.Errorf("expect %v, got %v", tt.expectTranslatedText, gotTranslatedText)
			}
		})
	}
}
