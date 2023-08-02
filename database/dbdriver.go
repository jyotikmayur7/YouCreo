package database

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/jyotikmayur7/YouCreo/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DatabaseClient(log hclog.Logger, ctx context.Context, config utils.Config) *mongo.Client {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.Database.URI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	log.Info("Connected to MongoDB!")

	return client
}
