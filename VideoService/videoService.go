package video_service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/go-hclog"
	"github.com/jyotikmayur7/YouCreo/api"
	"github.com/jyotikmayur7/YouCreo/database"
	"github.com/jyotikmayur7/YouCreo/models"
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

	videoModel := &models.Video{}
	videoTitle := req.GetVideoTitle()
	videoExtension := req.GetVideoExtension()
	videoDescription := req.GetVideoDescription()
	// need to add user ID before video title on the key to minatain unique key entry
	videoBlobReferenceKey := fmt.Sprintf("%s.%s", videoTitle, videoExtension)
	thumbnailBlobReferenceKey := videoTitle

	config := utils.GetConfig()
	ctx := stream.Context()

	videoData := bytes.Buffer{}
	videoThumbnail := bytes.Buffer{}

	var parts []*s3.CompletedPart
	partSize := int64(5 * 1024 * 1024)
	partNum := 1

	// How to handler error here? without return nil at the end of statement?
	go func() error {
		// Starting multipart upload process
		createdResp, errMP := vs.awsService.S3Client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
			Bucket: aws.String(config.Aws.Video.Bucket),
			Key:    aws.String(videoBlobReferenceKey),
		})
		if errMP != nil {
			vs.log.Error("Error creating multipart upload:", err)
			return errMP
		}

		for {
			vs.log.Info("Receiving video and thumbnail data")

			req, err := stream.Recv()
			if err == io.EOF {
				//Seding the last chunk of video
				partResp, err1 := vs.awsService.UploadPart(createdResp, videoData, partNum)
				if err1 != nil {
					return err1
				}
				parts = append(parts, &s3.CompletedPart{
					ETag:       partResp.ETag,
					PartNumber: aws.Int64(int64(partNum)),
				})
				partNum++

				// Uploading thumbnail
				_, err2 := vs.awsService.S3Client.PutObject(&s3.PutObjectInput{
					Body:   bytes.NewReader(videoThumbnail.Bytes()),
					Bucket: aws.String(config.Aws.Thumbnail.Bucket),
					Key:    aws.String(thumbnailBlobReferenceKey),
				})
				if err2 != nil {
					vs.log.Error("Failed to upload thumbnail ", err2)
					return err2
				}

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

		// Adding data in DB
		videoModel.Title = videoTitle
		videoModel.Description = videoDescription
		videoModel.PublishedTime = time.Now()
		videoModel.Likes = 0
		videoModel.Views = 0
		videoModel.VideoBlobReference = videoBlobReferenceKey
		videoModel.ThumbnailBlobReference = thumbnailBlobReferenceKey
		// Need to fetch channel name from user id or from channel table using user id
		// videoModel.ChannelName = ctx.Value("user_name")
		videoModel.ChannelName = "testChannel"

		dbErr := vs.DB.Video.CreateVideo(ctx, *videoModel)
		if dbErr != nil {
			vs.log.Error("Unable to add Video info on DB ", dbErr)
			return dbErr
		}

		return nil
	}()

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
