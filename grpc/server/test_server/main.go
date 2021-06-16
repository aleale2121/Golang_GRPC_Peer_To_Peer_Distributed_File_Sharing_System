package main

import (
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	fileServices "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/module/files/file_server_services"
	fileProto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/files"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/file_store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func main() {
	if err := RunServer(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func RunServer() error {

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get base path: %v", err)
	}
	store, err := file_store.NewStorage(basePath)
	if err != nil {
		log.Fatalf("cannot create storage: %v", err)
	}


	fs:=fileServices.NewGrpcFileServer(*store)
	gs := grpc.NewServer()



	fileProto.RegisterSongsServiceServer(gs,fs)

	reflection.Register(gs)

	grpcConnStr, err := constant.GetGrpcConnectionString()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", grpcConnStr)
	if err != nil {
		return err
	}
	err = gs.Serve(l)
	if err != nil {
		return err
	}
	fmt.Println("hello")
	return nil
}
