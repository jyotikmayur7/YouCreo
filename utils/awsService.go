package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/go-hclog"
)

type AWSService struct {
	S3Client *s3.Client
}

func NewAWSService(l hclog.Logger, ctx context.Context) *AWSService {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-1"))
	if err != nil {
		l.Error("Error while loading configurations", err)
		return nil
	}

	return &AWSService{
		S3Client: s3.NewFromConfig(cfg),
	}
}
