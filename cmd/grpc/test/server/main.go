package main

import (
	con "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/module/connection"
	protos "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/connect"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:7070")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	clientStore:=sync.Map{}
	liveServer:=con.NewLiveServer(&clientStore)

	// Start notifying  live  connections
	go liveServer.StartStreamLiveClients()

	// Register the server
	protos.RegisterLiveConnectionServer(grpcServer, liveServer)

	log.Printf("Starting server on address %s", lis.Addr().String())

	// Start listening
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
