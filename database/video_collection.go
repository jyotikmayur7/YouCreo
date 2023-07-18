package database

import (
	"github.com/jyotikmayur7/YouCreo/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoAccessor struct {
	Collection *mongo.Collection
}

type VideoCollection interface {
	CreateVideo(models.Video) (interface{}, error)
	GetAllVideos() ([]models.Video, error)
	DeleteVideoById(ID int) error
	UpdateVideo(models.Video) error
}
