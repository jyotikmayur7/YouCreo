package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseAccessor struct {
	Client *mongo.Client
	Video  VideoAccessor
}

func NewDatabaseAccessor(mc *mongo.Client) *DatabaseAccessor {
	return &DatabaseAccessor{Client: mc}
}

func (da *DatabaseAccessor) WithVideoAccessor(va VideoAccessor) *DatabaseAccessor {
	da.Video = va
	return da
}
