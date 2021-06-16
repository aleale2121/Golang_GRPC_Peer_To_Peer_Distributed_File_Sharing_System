package main

import (
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	connClientServices "github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/conn_client"
	fileServices "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/module/files/file_server_services"
	connectionProto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/connect"
	fileProto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/files"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/file_store"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strconv"
)
var ip string
var port int
func init() {
	port=GetPort()
	ip =fmt.Sprintf("127.0.0.1:%d",port)
}

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
	serverConn, err := CreateConnectionWithServer()
	if err != nil {
		log.Fatalf("cannot make connection: %v", err)
	}

	client, err := connClientServices.NewLiveClient(
		connectionProto.NewLiveConnectionClient(serverConn),
		serverConn,
		*GetInitialConnectionData(),
		)
	if err != nil {
		log.Fatal(err)
	}


	fileProto.RegisterSongsServiceServer(gs,fs)

	reflection.Register(gs)

	//starting the client
	go client.Start()

	l, err := net.Listen("tcp", ip)
	if err != nil {
		return err
	}
	err = gs.Serve(l)
	if err != nil {
		return err
	}
	return nil
}


func GetInitialConnectionData() *model.Info  {
	wd,_:=os.Getwd()
	pathX :=path.Join(wd,"assets","audio")
	files, err := ioutil.ReadDir(pathX)
	if err != nil {
		log.Fatal(err)
	}
	//validFormats:=[]string{
	//	"audio/basic","audio/L24",
	//	"audio/mid", "audio/mpeg","audio/mp4","audio/x-aiff","audio/x-mpegurl",
	//	"audio/ogg","audio/vorbis","audio/vnd.wav","text/plain",
	//}
	info:=model.Info{
		Id:    uuid.NewV4().String(),
		Port:  port,
		Files: make([]string,0),
	}
	for _, file := range files {
		if !file.IsDir() {
			info.Files= append(info.Files, file.Name())
		}
	}
	return &info
}

func GetPort() int {
	for i := 1; i < 65535; i++ {
		port := strconv.FormatInt(int64(i), 10)
		_, err := net.Dial("tcp", "127.0.0.1:" + port)
		if err != nil {
			fmt.Println("Port",i, "open")
			return  i
		}
	}
	return 0
}


func CreateConnectionWithServer() (*grpc.ClientConn, error) {
	return grpc.Dial("127.0.0.1:7070", []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}...)
}

