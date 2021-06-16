package main

import (
	"encoding/json"
	"fmt"
	protos "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/song"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/song_client_service"
	"log"
	"os"
)

func TestCreateSong(client *song_client_service.SongClient) {
	base, _ :=os.Getwd()
	audio:=base+"\\sample\\audio\\3.ogg"
	image:=base+"\\sample\\image\\1.png"

	resp, err := client.CreateSong(&protos.SongCreateInfo{
		ArtistId: "e0bc715f-a460-4af8-a346-0b9275fa5bfa",
		Title:    "Ethiopia",
		Duration: 56846456,
	},image,audio)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
	fmt.Println("-----------------------------------------------------------------")
}

func TestGetSong(client *song_client_service.SongClient) {
	resp,err:=client.GetSong("e0bc715f-a460-4af8-a346-0b9275fa5bfa")
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonX, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(jsonX))
	fmt.Println("--------------------------------------------------------------------")

}
func TestGetAllSongs(client *song_client_service.SongClient) {
	resp,err:=client.GetAllSongs("ASC","created_at",1,5)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonX, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(jsonX))
	fmt.Println("--------------------------------------------------------------------")
}

