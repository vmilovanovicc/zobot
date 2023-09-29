package backend

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"regexp"
	"testing"
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
	location := "http://amazon-polly-09232023111111.s3.amazonaws.com/"
	cases := []struct {
		client     func(t *testing.T) S3CreateBucketAPI
		bucketName string
		expect     string
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
					if e, a := "http://amazon-polly-09232023111111.s3.amazonaws.com/", location; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}

					return &s3.CreateBucketOutput{
						Location: &location,
					}, nil
				})
			},
			bucketName: "amazon-polly-09232023111111",
			expect:     "http://amazon-polly-09232023111111.s3.amazonaws.com/",
		},
	}

	for _, tt := range cases {
		t.Run("Test bucket creation", func(t *testing.T) {
			ctx := context.TODO()
			content, err := CreateBucketFromS3(ctx, tt.client(t), tt.bucketName)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := tt.expect, content; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
