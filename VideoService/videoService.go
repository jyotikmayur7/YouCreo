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
	DB         *database.DatabaseAccessor
	log        hclog.Logger
	awsService *utils.AWSService
}

func NewVideoService(l hclog.Logger, db *database.DatabaseAccessor, as *utils.AWSService) *VideoService {
	return &VideoService{
		DB:         db,
		log:        l,
		awsService: as,
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
	ctx := stream.Context()

	if err != nil {
		vs.log.Error("Error while loading configurations", err)
		return err
	}

	createdResp, err := vs.awsService.S3Client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(config.Aws.Video.Bucket),
		Key:    aws.String(videoTitle + "." + videoExtension),
	})
	if err != nil {
		vs.log.Error("Error creating multipart upload:", err)
		return err
	}

	videoBlobReference := createdResp.UploadId

	videoData := bytes.Buffer{}
	videoThumbnail := bytes.Buffer{}

	var parts []*s3.CompletedPart
	partSize := int64(5 * 1024 * 1024)
	partNum := 1

	go func() {
		for {
			vs.log.Info("Receiving video data")

			req, err := stream.Recv()
			if err == io.EOF {
				//Seding the last chunk of video
				partResp, err := vs.awsService.UploadPart(createdResp, videoData, partNum)
				if err != nil {
					return err
				}
				parts = append(parts, &s3.CompletedPart{
					ETag:       partResp.ETag,
					PartNumber: aws.Int64(int64(partNum)),
				})
				partNum++

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
				// Uploading the buffer chunk to s3
				partResp, err := vs.awsService.UploadPart(createdResp, videoData, partNum)
				if err != nil {
					return err
				}
				parts = append(parts, &s3.CompletedPart{
					ETag:       partResp.ETag,
					PartNumber: aws.Int64(int64(partNum)),
				})

				partNum++
			}
		}
		// Finishing the mulitpart part upload
		resp, err := vs.awsService.S3Client.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
			Bucket:   createdResp.Bucket,
			Key:      createdResp.Key,
			UploadId: createdResp.UploadId,
			MultipartUpload: &s3.CompletedMultipartUpload{
				Parts: parts,
			},
		})
		if err != nil {
			vs.log.Error("Error finishing mutipartupload: ", err)
		} else {
			vs.log.Info(resp.String())
		}
	}()

	//TODO: Upload thumbnail with single uplaod
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
