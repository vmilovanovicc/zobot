package backend

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	typesS3 "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"net/url"
	"strings"
	"time"
)

// CreateBucket creates an S3 bucket to store the output of Amazon Polly Synthesis Tasks i.e. audio streams.
func CreateBucket(bucketName string) (string, error) {
	cfg := LoadAWSConfig()
	client := s3.NewFromConfig(cfg)
	bucketName = fmt.Sprintf("%s-%s", "amazon-polly-audio", time.Now().Format("01022006150405"))

	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &typesS3.CreateBucketConfiguration{
			LocationConstraint: typesS3.BucketLocationConstraint(region),
		}, // basic configuration
	}
	ctx := context.TODO()
	resp, err := client.CreateBucket(ctx, params)
	if err != nil {
		log.Fatalf("failed to create s3 bucket, error:  %v\n", err)
		return "", err
	}
	fmt.Println(*resp.Location)
	return *resp.Location, nil
}

// GetBucketName retrieves the S3 bucket name.
func GetBucketName(locationURL string) (string, error) {
	parsedURL, err := url.Parse(locationURL)
	if err != nil {
		log.Fatalf("failed to parse S3 bucket locationURL, error: %v\n", err)
		return "", err
	}
	hostParts := strings.Split(parsedURL.Host, ".")
	if len(hostParts) < 3 {
		log.Fatalf("invalid s3 location URL: %s", locationURL)
	}
	outputS3BucketName := hostParts[0]
	fmt.Println(outputS3BucketName)
	return outputS3BucketName, nil
}

// GeneratePresignedURL retrieves a presigned URL for an Amazon S3 bucket object.
func GeneratePresignedURL(bucketName, objectName string, expires time.Time) (string, error) {
	cfg := LoadAWSConfig()
	client := s3.NewFromConfig(cfg)
	psClient := s3.NewPresignClient(client)

	parsedObjectName := strings.Split(objectName, "/")
	objectName = parsedObjectName[len(parsedObjectName)-1]
	expires = time.Now().Add(3600 * time.Second)

	params := &s3.GetObjectInput{
		Bucket:          aws.String(bucketName),
		Key:             aws.String(objectName),
		ResponseExpires: aws.Time(expires),
	}

	ctx := context.TODO()
	resp, err := psClient.PresignGetObject(ctx, params)
	if err != nil {
		log.Fatalf("failed to generate presigned URL: %v", err)
		return "", err
	}
	fmt.Println(resp.URL)
	return resp.URL, nil
}
