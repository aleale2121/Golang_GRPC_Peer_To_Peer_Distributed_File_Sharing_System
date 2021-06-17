package music

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	proto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/files"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/file_store"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MusicsClient struct {
	service proto.SongsServiceClient
	store   file_store.Storage
}

func NewMusicClient(rc *grpc.ClientConn, store file_store.Storage) *MusicsClient {
	return &MusicsClient{
		service: proto.NewSongsServiceClient(rc),
		store:   store,
	}
}

func (client *MusicsClient) UploadMusic(musicPath string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	stream, err := client.service.UploadSong(ctx)
	if err != nil {
		log.Println("cannot upload image: ", err)
		return "", err
	}

	musicAudio, err := os.Open(musicPath)
	if err != nil {
		log.Println("cannot open song: ", err)
		return "", err
	}
	pathSplitted := strings.Split(musicAudio.Name(), "/")
	fmt.Println("audio name --")
	defer musicAudio.Close()

	req := &proto.UploadSongRequest{
		Data: &proto.UploadSongRequest_Title{
			Title: pathSplitted[len(pathSplitted)-1]},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send song info to server: ", err, stream.RecvMsg(nil))
	}

	readerSong := bufio.NewReader(musicAudio)
	musicBuffers := make([]byte, 1024)

	for {
		n, err := readerSong.Read(musicBuffers)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to song buffers: ", err)
		}

		req := &proto.UploadSongRequest{
			Data: &proto.UploadSongRequest_ChunkData{ChunkData: musicBuffers[:n]},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send song chunk to server: ", err, stream.RecvMsg(nil))
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	return res.Id, err
}

func (client *MusicsClient) DownloadMusic(fileId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	req := &proto.DownloadSongRequest{
		SongId: fileId,
	}
	stream, err := client.service.DownloadSong(ctx, req)
	if err != nil {
		log.Fatal("cannot download image: ", err)
	}
	defer stream.CloseSend()
	buffer := bytes.Buffer{}
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			// we've reached the end of stream
			log.Println("recived all chunks")
			break
		}
		if err != nil {
			log.Fatalf("error while reciving chunk %v", err)
		}
		chunk := msg.GetChunkData()

		_, err = buffer.Write(chunk)
		if err != nil {
			log.Fatalf("couldn't write chunk data: %v", err)
		}
	}

	err = client.saveFile(fileId, "audio", buffer)
	if err != nil {
		return fmt.Errorf("cannot write song cover to file: %w", err)
	}

	return nil
}

func (client *MusicsClient) saveFile(id, path string, buffer bytes.Buffer) error {

	fp := filepath.Join("assets", path, id)
	err := client.store.SaveChunk(fp, buffer)
	if err != nil {
		return err
	}
	return nil
}
