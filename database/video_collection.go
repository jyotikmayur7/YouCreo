package database

import (
	"context"

	"github.com/jyotikmayur7/YouCreo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoAccessor struct {
	Collection *mongo.Collection
}

func NewVideoAccessor(col *mongo.Collection) *VideoAccessor {
	return &VideoAccessor{Collection: col}
}

type VideoCollection interface {
	CreateVideo(ctx context.Context, video models.Video) error
	GetAllVideos(ctx context.Context) ([]models.Video, error)
	GetVideoById(ctx context.Context, ID primitive.ObjectID) (models.Video, error)
	DeleteVideoById(ctx context.Context, ID primitive.ObjectID) error
	UpdateVideo(ctx context.Context, video models.Video) error
}

func (v *VideoAccessor) CreateVideo(ctx context.Context, video models.Video) error {
	insert := video.ToBson()

	_, err := v.Collection.InsertOne(ctx, insert)
	if err != nil {
		return err
	}

	return nil
}

func (v *VideoAccessor) GetAllVideos(ctx context.Context) ([]models.Video, error) {
	var allVideos []models.Video

	cursor, err := v.Collection.Find(ctx, bson.D{})
	if err != nil {
		return allVideos, err
	}

	for cursor.Next(ctx) {
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

	err = cursor.Close(ctx)

	return allVideos, err

}

func (v *VideoAccessor) GetVideoById(ctx context.Context, ID primitive.ObjectID) (models.Video, error) {
	var video models.Video

	filter := bson.D{primitive.E{Key: "_id", Value: ID}}

	err := v.Collection.FindOne(ctx, filter).Decode(&video)
	if err != nil {
		return video, err
	}

	return video, err
}

func (v *VideoAccessor) DeleteVideoById(ctx context.Context, ID primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	_, err := v.Collection.DeleteOne(ctx, filter)

	return err

}

func (v *VideoAccessor) UpdateVideo(ctx context.Context, video models.Video) error {
	filter := bson.D{primitive.E{Key: "_id", Value: video.ID}}
	update := video.ToBson()
	_, err := v.Collection.UpdateOne(ctx, filter, update)

	return err
}
