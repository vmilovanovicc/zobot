package backend

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"regexp"
	"testing"
	"time"
)

type mockCreateBucketAPI func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)

func (m mockCreateBucketAPI) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return m(ctx, params, optFns...)
}

// isValidBucketName is a helper function to check the S3 bucket naming rules.
// - Must begin with a letter or number
// - Must end with a letter or number
// - Can consist only of lowercase letters, numbers, dots (.), and hyphens (-)
func isValidBucketName(bucketName string) bool {
	pattern := "^[a-zA-Z0-9][a-z0-9.-]*[a-z0-9]$"
	match, _ := regexp.MatchString(pattern, bucketName)
	return match
}

func TestCreateBucketFromS3(t *testing.T) {
	cases := []struct {
		client         func(t *testing.T) S3CreateBucketAPI
		bucketName     string
		expectLocation string
	}{
		{
			client: func(t *testing.T) S3CreateBucketAPI {
				return mockCreateBucketAPI(func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
					t.Helper()
					if params.Bucket == nil {
						t.Fatalf("expect bucket name to not be nil")
					}
					if len(*params.Bucket) < 3 || len(*params.Bucket) > 63 {
						t.Fatalf("expect bucket name to be between 3 and 63 characters long")
					}
					if !isValidBucketName(*params.Bucket) {
						t.Fatalf("expect bucket to meet S3 bucket naming rules")
					}

					return &s3.CreateBucketOutput{
						Location: aws.String("http://amazon-polly-09232023111111.s3.amazonaws.com/"),
					}, nil
				})
			},
			bucketName:     "amazon-polly-09232023111111",
			expectLocation: "http://amazon-polly-09232023111111.s3.amazonaws.com/",
		},
	}

	for _, tt := range cases {
		t.Run("Test bucket creation", func(t *testing.T) {
			ctx := context.TODO()
			gotLocation, err := CreateBucketFromS3(ctx, tt.client(t), tt.bucketName)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if gotLocation != tt.expectLocation {
				t.Errorf("expect %v, got %v", tt.expectLocation, gotLocation)

			}
		})
	}
}

func TestGetBucketName(t *testing.T) {
	type args struct {
		locationURL string
	}
	tests := []struct {
		name                   string
		args                   args
		wantOutputS3BucketName string
		wantErr                bool
	}{
		{
			name: "Test retrieving S3 bucket name",
			args: args{
				locationURL: "http://amazon-polly-09232023111111.s3.amazonaws.com/",
			},
			wantOutputS3BucketName: "amazon-polly-09232023111111",
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutputS3BucketName, err := GetBucketName(tt.args.locationURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutputS3BucketName != tt.wantOutputS3BucketName {
				t.Errorf("got OutputS3BucketName %v, want %v", gotOutputS3BucketName, tt.wantOutputS3BucketName)
			}
		})
	}
}

type mockGetPresignedURLAPI func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)

func (m mockGetPresignedURLAPI) GetPresignedURL(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	return m(ctx, params, optFns...)
}

func TestGetPresignedURLFromS3(t *testing.T) {
	cases := []struct {
		client     func(t *testing.T) S3PresignGetObjectAPI
		bucketName string
		objectName string
		expires    time.Time
		expectURL  string
	}{
		{
			client: func(t *testing.T) S3PresignGetObjectAPI {
				return mockGetPresignedURLAPI(func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
					t.Helper()

					if *params.Bucket == "" {
						t.Fatalf("expect bucket name to not be nil")
					}

					if *params.Key == "" {
						t.Fatalf("expect object name to not be nil")
					}

					if params.ResponseExpires.IsZero() {
						t.Fatalf("expect time to be set, currently not set")
					}

					return &v4.PresignedHTTPRequest{
						URL: "https://mock-bucket.s3.eu-central-1.amazonaws.com/mock-object.mp3?XXXX",
					}, nil
				})
			}, // end client
			bucketName: "mock-bucket",
			objectName: "mock-object",
			expires:    time.Now().Add(3600 * time.Second),
			expectURL:  "https://mock-bucket.s3.eu-central-1.amazonaws.com/mock-object.mp3?XXXX",
		},
	}
	for _, tt := range cases {
		t.Run("Test presigned URL", func(t *testing.T) {
			ctx := context.TODO()

			gotPresignedURL, err := GetPresignedURLFromS3(ctx, tt.client(t), tt.bucketName, tt.objectName, tt.expires)
			if err != nil {
				t.Fatalf("expect no error, got: %v", err)
			}
			if tt.expectURL != gotPresignedURL {
				t.Errorf("expect %v, got %v", tt.expectURL, gotPresignedURL)
			}
		})
	}
}
