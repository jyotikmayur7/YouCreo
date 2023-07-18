package database

import "go.mongodb.org/mongo-driver/mongo"

type DatabaseAccessor struct {
	Client *mongo.Client
}
