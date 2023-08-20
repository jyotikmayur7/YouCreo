package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/go-hclog"
	video_service "github.com/jyotikmayur7/YouCreo/VideoService"
	"github.com/jyotikmayur7/YouCreo/api"
	"github.com/jyotikmayur7/YouCreo/database"
	"github.com/jyotikmayur7/YouCreo/middleware"
	"github.com/jyotikmayur7/YouCreo/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func StartService() {
	ctx := context.Background()
	log := hclog.Default()
	config, err := utils.LoadConfig(log)
	if err != nil {
		log.Error(err.Error())
	}
	databaseAccessor := database.NewDatabaseAccessor(database.DatabaseClient(log, ctx))
	databaseAccessor = initDatabaseAccessor(databaseAccessor, config)
	awsService := utils.NewAWSService(log)

	grpcServer := grpc.NewServer()
	videoService := video_service.NewVideoService(log, databaseAccessor, awsService)

	api.RegisterVideoServiceServer(grpcServer, videoService)

	reflection.Register(grpcServer)

	// Healthend point is required

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Server.Grpc.Port))
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	log.Info("Serving gRPC on ", fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Grpc.Port))
	go func() {
		err := grpcServer.Serve(l)
		if err != nil {
			log.Error(err.Error())
		}
	}()

	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Grpc.Port),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to dial server:", err)
	}

	gatewayMux := runtime.NewServeMux()

	err = api.RegisterVideoServiceHandler(ctx, gatewayMux, conn)
	if err != nil {
		log.Error("Failed to register gateway:", err)
	}

	gatewayServer := &http.Server{
		Addr:    ":" + config.Server.Gateway.Port,
		Handler: middleware.AddContext(ctx, gatewayMux),
	}

	log.Info("Serving gRPC-Gateway on ", fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Gateway.Port))
	go func() {
		err := gatewayServer.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Error("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(ctx, 30*time.Second)
	gatewayServer.Shutdown(tc)
	grpcServer.GracefulStop()
	err = databaseAccessor.Client.Disconnect(ctx)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Servers gracefully stopped!")
}

func initDatabaseAccessor(da *database.DatabaseAccessor, config *utils.Config) *database.DatabaseAccessor {
	db := da.Client.Database(config.Database.Name)
	videoCollection := db.Collection(config.Database.Collection.Video)
	videoAccessor := database.NewVideoAccessor(videoCollection)
	da = da.WithVideoAccessor(*videoAccessor)
	return da
}
