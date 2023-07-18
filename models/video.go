package models

type Video struct {
	ID            int       `bson:"_id" json:"video_id"`
	Title         string    `bson:"title" json:"video_title"`
	Description   string    `bson:"description" json:"video_description"`
	PublishedTime time.time `bson:"published_time" json:"published_time"`
	PublishedDate time.date `bson:"published_date" json:"published_date"`
	Likes         int64     `bson:"likes" json:"likes"`
	Views         int64     `bson:"views" json:"views"`
	ChannelName   string    `bson:"channel_name" json:"channel_name"`
}
