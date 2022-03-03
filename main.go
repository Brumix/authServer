package main

import (
	_ "authServer/config"
	"authServer/core"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

func init() {
	log.Info("GRPC AUTH SERVER!!!")
}

func main() {

	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPCPORT"))
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	var grpcServer = grpc.NewServer()

	core.RegisterAuthServiceServer(grpcServer, &core.Server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to to server the grpc Server: %v", err)
	}
}
