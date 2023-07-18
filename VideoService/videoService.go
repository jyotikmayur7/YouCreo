package video_service

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/jyotikmayur7/YouCreo/api"
	"github.com/jyotikmayur7/YouCreo/database"
)

type VideoService struct {
	DB  *database.DatabaseAccessor
	log hclog.Logger
	api.UnimplementedVideoServiceServer
}

func NewVideoService(l hclog.Logger, db *database.DatabaseAccessor) *VideoService {
	return &VideoService{
		DB:                              db,
		log:                             l,
		UnimplementedVideoServiceServer: api.UnimplementedVideoServiceServer{},
	}
}
func (vs *VideoService) CreateVideo(stream api.VideoService_CreateVideoServer) error {
	return nil
}

func (vs *VideoService) DeleteVideo(context.Context, *api.DeleteVideoRequest) (*api.DeleteVideoResponse, error) {
	return nil, nil
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
