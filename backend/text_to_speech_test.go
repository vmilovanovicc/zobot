package backend

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	typesPolly "github.com/aws/aws-sdk-go-v2/service/polly/types"
	"testing"
)

func TestGetTargetVoice(t *testing.T) {
	type args struct {
		language string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "English - Successful Match",
			args: args{
				language: "en",
			},
			want:    "Emma",
			wantErr: false,
		},
		{
			name: "Not supported language - Want Err",
			args: args{
				language: "xxxxxxxx",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTargetVoice(tt.args.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf(" got error %v, want %v", got, tt.want)
			}
		})
	}
}

type mockStartSpeechSynthesisTaskAPI func(ctx context.Context, params *polly.StartSpeechSynthesisTaskInput, optFns ...func(options *polly.Options)) (*polly.StartSpeechSynthesisTaskOutput, error)

func (m mockStartSpeechSynthesisTaskAPI) StartSpeechSynthesisTask(ctx context.Context, params *polly.StartSpeechSynthesisTaskInput, optFns ...func(options *polly.Options)) (*polly.StartSpeechSynthesisTaskOutput, error) {
	return m(ctx, params, optFns...)
}

func TestStartSpeechSynthesisTaskFromPolly(t *testing.T) {
	cases := []struct {
		client           func(t *testing.T) PollyStartSpeechSynthesisTaskAPI
		text             string
		bucketName       string
		languageCode     string
		targetVoice      string
		expectTaskId     string
		expectObjectName string
	}{
		{
			client: func(t *testing.T) PollyStartSpeechSynthesisTaskAPI {
				return mockStartSpeechSynthesisTaskAPI(func(ctx context.Context, params *polly.StartSpeechSynthesisTaskInput, optFns ...func(options *polly.Options)) (*polly.StartSpeechSynthesisTaskOutput, error) {
					t.Helper()
					if *params.Text == "" {
						t.Fatalf("expect text to not be nil")
					}
					if *params.OutputS3BucketName == "" {
						t.Fatalf("expect bucket name to not be nil")
					}
					if params.LanguageCode == "" {
						t.Fatalf("expect language name to not be nil ")
					}
					if params.VoiceId == "" {
						t.Fatalf("expect target voice not to be nil")
					}

					return &polly.StartSpeechSynthesisTaskOutput{
						SynthesisTask: &typesPolly.SynthesisTask{
							TaskId:    aws.String("mock-task-id"),
							OutputUri: aws.String("mock-task-id.mp3"),
						},
					}, nil
				})
			},
			text:             "Hello",
			bucketName:       "amazon-polly-audio",
			languageCode:     "es",
			targetVoice:      "Lucia",
			expectTaskId:     "mock-task-id",
			expectObjectName: "mock-task-id.mp3",
		},
	}

	for _, tt := range cases {
		t.Run("Test speech synthesis", func(t *testing.T) {
			ctx := context.TODO()
			gotTaskId, gotObjectName, err := StartSpeechSynthesisTaskFromPolly(ctx, tt.client(t), tt.text, tt.bucketName, tt.languageCode, tt.targetVoice)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if gotTaskId != tt.expectTaskId {
				t.Errorf("got taskId %v, want %v", gotTaskId, tt.expectTaskId)
			}
			if gotObjectName != tt.expectObjectName {
				t.Errorf("got objectName %v, want %v", gotObjectName, gotTaskId)
			}
		})
	}
}

type mockGetSpeechSynthesisTaskAPI func(ctx context.Context, params *polly.GetSpeechSynthesisTaskInput, optFns ...func(*polly.Options)) (*polly.GetSpeechSynthesisTaskOutput, error)

func (m mockGetSpeechSynthesisTaskAPI) GetSpeechSynthesisTask(ctx context.Context, params *polly.GetSpeechSynthesisTaskInput, optFns ...func(*polly.Options)) (*polly.GetSpeechSynthesisTaskOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetSpeechSynthesisTaskFromPolly(t *testing.T) {
	cases := []struct {
		client           func(t *testing.T) PollyGetSpeechSynthesisTaskAPI
		taskId           string
		expectTaskStatus typesPolly.TaskStatus
	}{
		{
			client: func(t *testing.T) PollyGetSpeechSynthesisTaskAPI {
				return mockGetSpeechSynthesisTaskAPI(func(ctx context.Context, params *polly.GetSpeechSynthesisTaskInput, optFns ...func(*polly.Options)) (*polly.GetSpeechSynthesisTaskOutput, error) {
					t.Helper()
					if *params.TaskId == "" {
						t.Fatalf("expect task id to not be nil")
					}

					return &polly.GetSpeechSynthesisTaskOutput{
						SynthesisTask: &typesPolly.SynthesisTask{
							TaskStatus: "mock-status-scheduled",
						},
					}, nil
				})
			},
			taskId:           "mock-task-id-123",
			expectTaskStatus: "mock-status-scheduled",
		},
	}
	for _, tt := range cases {
		t.Run("Test task status", func(t *testing.T) {
			ctx := context.TODO()
			gotTaskStatus, err := GetSpeechSynthesisTaskFromPolly(ctx, tt.client(t), tt.taskId)
			if err != nil {
				t.Fatalf("expect no error, got: %v", err)
			}
			if tt.expectTaskStatus != gotTaskStatus {
				t.Errorf("expect %v, got %v", tt.expectTaskStatus, gotTaskStatus)
			}
		})
	}
}
