package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID                     primitive.ObjectID `bson:"_id" json:"video_id"`
	Title                  string             `bson:"title" json:"video_title"`
	Description            string             `bson:"description" json:"video_description"`
	PublishedTime          time.Time          `bson:"published_time" json:"published_time"`
	Likes                  int64              `bson:"likes" json:"likes"`
	Views                  int64              `bson:"views" json:"views"`
	ChannelName            string             `bson:"channel_name" json:"channel_name"`
	VideoBlobReference     string             `bson:"video_blob_reference" json:"video_blob_reference"`
	ThumbnailBlobReference string             `bson:"thumbnail_blob_reference" json:"thumbnail_blob_reference"`
}

func (v *Video) ToBson() bson.M {
	return bson.M{
		"title":                    v.Title,
		"description":              v.Description,
		"published_time":           v.PublishedTime,
		"likes":                    v.Likes,
		"views":                    v.Views,
		"channel_name":             v.ChannelName,
		"video_blob_reference":     v.VideoBlobReference,
		"thumbnail_blob_reference": v.ThumbnailBlobReference,
	}
}
