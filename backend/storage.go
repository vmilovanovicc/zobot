package backend

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"net/url"
	"strings"
	"time"
)

// CreateBucket creates an S3 bucket to store the output of Amazon Polly Synthesis Tasks i.e. audio streams.
func CreateBucket(bucketName string) (location string, err error) {
	cfg := LoadAWSConfig()
	client := s3.NewFromConfig(cfg)
	bucketName = fmt.Sprintf("%s-%s", "amazon-polly-audio", time.Now().Format("01022006150405"))

	params := &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
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
