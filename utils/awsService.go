package utils

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/go-hclog"
)

type AWSService struct {
	S3Client *s3.S3
	log      hclog.Logger
}

func NewAWSService(l hclog.Logger) *AWSService {
	sysConfig := GetConfig()

	return &AWSService{
		S3Client: s3.New(session.Must(
			session.NewSession(&aws.Config{
				Region: aws.String(sysConfig.Aws.Region),
			}))),
		log: l,
	}
}

func (awsServ *AWSService) UploadPart(createdResp *s3.CreateMultipartUploadOutput, videoData bytes.Buffer, partNum int) (completedPart *s3.CompletedPart, err error) {
	partResp, err := awsServ.S3Client.UploadPart(&s3.UploadPartInput{
		Body:          bytes.NewReader(videoData.Bytes()),
		Bucket:        createdResp.Bucket,
		Key:           createdResp.Bucket,
		UploadId:      createdResp.UploadId,
		ContentLength: aws.Int64(int64(videoData.Len())),
		PartNumber:    aws.Int64(int64(partNum)),
	})
	if err != nil {
		awsServ.log.Error("Error uploading part: ", err)
		awsServ.log.Info("Aborting the multipart upload for partnumber: ", partNum)
		_, err = awsServ.S3Client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
			Bucket:   createdResp.Bucket,
			Key:      createdResp.Key,
			UploadId: createdResp.UploadId,
		})
		if err != nil {
			awsServ.log.Error("Error aborting upload for part number: ", partNum, " ", err)
		}
		return nil, err
	}

	return &s3.CompletedPart{
		ETag:       partResp.ETag,
		PartNumber: aws.Int64(int64(partNum)),
	}, nil
}
