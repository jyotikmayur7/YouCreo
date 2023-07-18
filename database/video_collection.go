package database

import "go.mongodb.org/mongo-driver/mongo"

type VideoAccessor struct {
	Collection *mongo.Collection
}
