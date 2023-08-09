package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSService struct {
	S3Client *s3.Client
}

func NewAWSService(ctx context.Context) (*AWSService, error) {
	sysConfig := GetConfig()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(sysConfig.Aws.Region))
	if err != nil {
		return nil, err
	}

	return &AWSService{
		S3Client: s3.NewFromConfig(cfg),
	}, nil
}
