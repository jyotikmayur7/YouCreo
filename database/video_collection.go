package database

import (
	"context"

	"github.com/jyotikmayur7/YouCreo/models"
	"go.mongodb.org/mongo-driver/bson"
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
	var allVideos []models.Video

	cursor, err := v.Collection.Find(v.ctx, bson.D{})
	if err != nil {
		return allVideos, err
	}

	for cursor.Next(v.ctx) {
		var video models.Video
		err := cursor.Decode(&video)
		if err != nil {
			return allVideos, err
		}
		allVideos = append(allVideos, video)
	}

	if err := cursor.Err(); err != nil {
		return allVideos, nil
	}

	err = cursor.Close(v.ctx)

	return allVideos, err

}

func (v *VideoAccessor) DeleteVideoById(ID int) error {
	return nil
}

func (v *VideoAccessor) UpdateVideo(video models.Video) error {
	return nil
}
