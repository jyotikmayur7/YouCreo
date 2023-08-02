package video_service

import (
	"bytes"
	"context"
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/jyotikmayur7/YouCreo/api"
	"github.com/jyotikmayur7/YouCreo/database"
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
	videoDescription := req.GetVideoDescription()

	videoData := bytes.Buffer{}
	videoThumbnail := bytes.Buffer{}

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
		// Instead of storing data locally need to send this chunk to aws S3 to store it so that bytes.Buffer{} won't exceed the size limit
		videoData.Write(videoChunk)
	}

	//TODO's

	// AWS s3 store video call
	//blobReference :=

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
