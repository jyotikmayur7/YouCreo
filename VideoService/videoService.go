package video_service

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/go-hclog"
	"github.com/jyotikmayur7/YouCreo/api"
	"github.com/jyotikmayur7/YouCreo/database"
	"github.com/jyotikmayur7/YouCreo/utils"
)

type VideoService struct {
	DB  *database.DatabaseAccessor
	log hclog.Logger
}

func NewVideoService(l hclog.Logger, db *database.DatabaseAccessor) *VideoService {
	return &VideoService{
		DB:  db,
		log: l,
	}
}

func (vs *VideoService) CreateVideo(stream api.VideoService_CreateVideoServer) error {
	req, err := stream.Recv()
	if err != nil {
		vs.log.Error("Error: Unable to receive video info ", err)
		return err
	}

	videoTitle := req.GetVideoTitle()
	videoExtension := req.GetVideoExtension()
	videoDescription := req.GetVideoDescription()

	config := utils.GetConfig()
	ctx := utils.GetContext()
	awsService, err := utils.NewAWSService(ctx)
	if err != nil {
		vs.log.Error("Error while loading configurations", err)
		return err
	}

	createMultipartUploadInput := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(config.Aws.Video.Bucket),
		Key:    aws.String(videoTitle + "." + videoExtension),
	}

	createdResp, err := awsService.S3Client.CreateMultipartUpload(ctx, createMultipartUploadInput)
	if err != nil {
		vs.log.Error("Error creating multipart upload:", err)
		return err
	}

	videoBlobReference := createdResp.UploadId

	partSize := int64(5 * 1024 * 1024)

	//var videoBuffer []byte = make([]byte, partSize)
	videoData := bytes.Buffer{}
	videoThumbnail := bytes.Buffer{}

	var parts []*s3.CompletedPart
	for {
		vs.log.Info("Receiving video data")

		req, err := stream.Recv()
		if err == io.EOF {
			vs.log.Info("Client has stopped the upload, upload is finished ", err)
			break
		}
		if err != nil {
			vs.log.Error("Error: Unable to read video chunk from client ", err)
			return err
		}

		thumbnailChunk := req.GetVideoThumbnail()
		if thumbnailChunk != nil {
			videoThumbnail.Write(thumbnailChunk)
		}

		videoChunk := req.GetVideoContent()
		videoData.Write(videoChunk)
		if int64(videoData.Len()) >= partSize {
			// Upload the buffer to s3
			partInput := &s3.UploadPartInput{
				Body:          bytes.NewReader(videoData.Bytes()),
				Bucket:        aws.String(config.Aws.Video.Bucket),
				Key:           aws.String(videoTitle + "." + videoExtension),
				UploadId:      videoBlobReference,
				ContentLength: *aws.Int64(int64(videoData.Len())),
				PartNumber:    *aws.Int32(int32(len(parts)) + 1),
			}

			partResp, err := awsService.S3Client.UploadPart(ctx, partInput)
			if err != nil {
				vs.log.Error("Error uploading part: ", err)
				return err
			}

			parts = append(parts, &s3.UploadPartOutput{
				ETag: partResp.ETag,
			})
		}
	}

	//TODO after multipart upload:
	//Store the info on database

	return nil
}

func (vs *VideoService) DeleteVideo(context.Context, *api.DeleteVideoRequest) (*api.DeleteVideoResponse, error) {
	vs.log.Info("Delete Method")
	return &api.DeleteVideoResponse{}, nil
}
func (vs *VideoService) SteamVideo(req *api.StreamVideoRequest, stream api.VideoService_SteamVideoServer) error {
	return nil
}
func (vs *VideoService) UpdateVideo(stream api.VideoService_UpdateVideoServer) error {
	return nil
}
func (vs *VideoService) GetAllVideos(req *api.GetAllVideosRequest, stream api.VideoService_GetAllVideosServer) error {
	return nil
}
