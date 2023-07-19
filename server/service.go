package server

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/hashicorp/go-hclog"
	video_service "github.com/jyotikmayur7/YouCreo/VideoService"
	"github.com/jyotikmayur7/YouCreo/api"
	"github.com/jyotikmayur7/YouCreo/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartService() {
	ctx := context.Background()
	log := hclog.Default()
	grpcServer := grpc.NewServer()
	databaseAccessor := database.NewDatabaseAccessor(database.DatabaseClient(log, ctx))
	databaseAccessor = initDatabaseAccessor(databaseAccessor)
	videoService := video_service.NewVideoService(log, databaseAccessor)

	api.RegisterVideoServiceServer(grpcServer, videoService)

	reflection.Register(grpcServer)

	// http server is required
	// 443 -> http, 8000 -> grpc
	// Healthend point is required

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	go func() {
		err := grpcServer.Serve(l)
		if err != nil {
			log.Error(err.Error())
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Error("Received terminate, graceful shutdown", sig)

	grpcServer.GracefulStop()
	err = databaseAccessor.Client.Disconnect(ctx)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Server gracefully stopped!")
}

func initDatabaseAccessor(da *database.DatabaseAccessor) *database.DatabaseAccessor {
	db := da.Client.Database("YouCreo")
	videoCollection := db.Collection("Video")
	videoAccessor := database.NewVideoAccessor(videoCollection)
	da.WithVideoAccessor(*videoAccessor)
	return da
}
