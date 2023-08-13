package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSService struct {
	S3Client *s3.S3
}

func NewAWSService() *AWSService {
	sysConfig := GetConfig()

	return &AWSService{
		S3Client: s3.New(session.Must(
			session.NewSession(&aws.Config{
				Region: aws.String(sysConfig.Aws.Region),
			}))),
	}
}
