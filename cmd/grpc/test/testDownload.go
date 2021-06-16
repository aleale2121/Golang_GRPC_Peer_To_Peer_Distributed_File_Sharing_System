package main

import (
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/music_client_service"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/file_store"
	"google.golang.org/grpc"
	"log"
	"os"
)

func main() {

	transportOption := grpc.WithInsecure()
	cc2, err := grpc.Dial(
		"127.0.0.1:2",
		transportOption,
	)
	if err != nil {
		log.Println("cannot dial server: ", err)
	}
	basePath, err := os.Getwd()
	if err != nil {
		log.Printf("cannot get base path: %v\n", err)
	}
	store, err := file_store.NewStorage(basePath)
	if err != nil {
		log.Printf("cannot create storage: %v\n", err)
	}

	musicClient := music_client_service.NewMusicClient(cc2, *store)
	err = musicClient.DownloadFile("3.ogg")
	if err != nil {
		fmt.Println("Error Occurred while downloading music")
		fmt.Println(err)
		return
	}
}
