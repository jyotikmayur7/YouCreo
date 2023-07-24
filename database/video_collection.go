package database

import (
	"context"

	"github.com/jyotikmayur7/YouCreo/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoAccessor struct {
	Collection *mongo.Collection
	ctx        context.Context
}

func NewVideoAccessor(col *mongo.Collection, ctx context.Context) *VideoAccessor {
	return &VideoAccessor{Collection: col}
}

type VideoCollection interface {
	CreateVideo(models.Video) error
	GetAllVideos() ([]models.Video, error)
	DeleteVideoById(ID int) error
	UpdateVideo(models.Video) error
}

func (v *VideoAccessor) CreateVideo(video models.Video) error {
	insert := video.ToBson()

	_, err := v.Collection.InsertOne(v.ctx, insert)
	if err != nil {
		return err
	}

	return nil
}

func (v *VideoAccessor) GetAllVideos() ([]models.Video, error) {
	return nil, nil
}

func (v *VideoAccessor) DeleteVideoById(ID int) error {
	return nil
}

func (v *VideoAccessor) UpdateVideo(video models.Video) error {
	return nil
}
