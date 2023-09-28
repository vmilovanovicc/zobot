package backend

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
)

var region = "eu-central-1"

// LoadAWSConfig loads the Shared AWS Configuration (~/.aws/config)
func LoadAWSConfig() (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load configuration, %v\n", err)
	}
	return cfg
}
