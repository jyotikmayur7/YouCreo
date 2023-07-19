package database

import (
	"github.com/jyotikmayur7/YouCreo/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoAccessor struct {
	Collection *mongo.Collection
}

func NewVideoAccessor(col *mongo.Collection) *VideoAccessor {
	return &VideoAccessor{Collection: col}
}

type VideoCollection interface {
	CreateVideo(models.Video) error
	GetAllVideos() ([]models.Video, error)
	DeleteVideoById(ID int) error
	UpdateVideo(models.Video) error
}

func (v *VideoAccessor) CreateVideo(video models.Video) error {
	return nil
}
