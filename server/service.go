package server

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	video_service "github.com/jyotikmayur7/YouCreo/VideoService"
	"github.com/jyotikmayur7/YouCreo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartService() {
	log := hclog.Default()
	grpcServer := grpc.NewServer()
	videoService := video_service.NewVideoService(log)

	api.RegisterVideoServiceServer(grpcServer, videoService)

	reflection.Register(grpcServer)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	grpcServer.Serve(l)

	// Go routine code required for this to work
	// sigChan := make(chan os.Signal)
	// signal.Notify(sigChan, os.Interrupt)
	// signal.Notify(sigChan, os.Kill)

	// sig := <-sigChan
	// l.Println("Received terminate, graceful shutdown", sig)

	// tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// s.Shutdown(tc)
}
