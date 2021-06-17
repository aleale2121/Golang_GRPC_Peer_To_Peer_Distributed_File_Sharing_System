package files

import (
	"bytes"
	"context"
	"fmt"
	proto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/files"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/file_store"
	"github.com/gabriel-vasile/mimetype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	songs = make([]Info, 0)
)

type Info struct {
	id    string
	port  int
	songs []string
}
type grpcFileServiceServer struct {
	store file_store.Storage
	proto.UnimplementedSongsServiceServer
}

func (s grpcFileServiceServer) Connect(ctx context.Context, request *proto.ConnectRequest) (*proto.ConnectResponse, error) {
	songs = append(songs, *ConvertProtoToSong(request.Info))
	return &proto.ConnectResponse{}, nil
}
func NewGrpcFileServer(store file_store.Storage) proto.SongsServiceServer {
	return &grpcFileServiceServer{
		store: store,
	}
}
func (s grpcFileServiceServer) UploadSong(stream proto.SongsService_UploadSongServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive image info")
	}

	title := req.GetTitle()

	if title == "" {
		return status.Errorf(codes.InvalidArgument, "Song Title Cannot Be Empty")
	}
	buffer := bytes.Buffer{}
	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err)
		}

		chunkData := req.GetChunkData()

		_, err = buffer.Write(chunkData)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write  song chunk data: %v", err)
		}

	}
	mime := mimetype.Detect(buffer.Bytes())
	//if !mimetype.EqualsAny(mime.String(),"image/jpeg","image/pjpeg",
	//	"image/png", "image/tiff","image/x-tiff","image/vnd.wap.wbmp"){
	//	return status.Errorf(codes.InvalidArgument, "the cover image you upload is not image")
	//}

	songId := title + mime.Extension()
	err = s.saveFile(songId, "audio", buffer)
	if err != nil {
		return fmt.Errorf("cannot write song cover to file: %w", err)
	}

	res := &proto.UploadSongResponse{
		Id: title,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}
	log.Println("song saved")
	return nil
}

func (s grpcFileServiceServer) GetSongsList(ctx context.Context, request *proto.GetSongsRequest) (*proto.GetSongsResponse, error) {
	songsInfo := make([]*proto.SongData, 0)
	for _, sg := range songs {
		songsInfo = append(songsInfo, ConvertToProtoSong(&sg))
	}
	return &proto.GetSongsResponse{Songs: songsInfo}, nil
}

func (s grpcFileServiceServer) DownloadSong(request *proto.DownloadSongRequest, server proto.SongsService_DownloadSongServer) error {
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "assets", "audio", request.SongId)

	fileX, err := os.Open(fp)
	if err != nil {
		log.Fatal("cannot open coverImage: ", err)
	}
	defer fileX.Close()
	fmt.Println("------fileX opened---", fileX.Name())
	buff := make([]byte, 1024)
	for {
		bytesRead, err := fileX.Read(buff)
		if err == io.EOF {
			log.Println("End of file")
			break
		} else if err != nil {
			log.Println("error--", err)
			break
		}
		resp := &proto.DownloadSongResponse{
			ChunkData: buff[:bytesRead],
		}
		err = server.Send(resp)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}

	}
	return nil
}

func ConvertProtoToSong(song *proto.SongData) *Info {
	return &Info{
		id:    song.Id,
		port:  int(song.Port),
		songs: song.Songs,
	}
}

func ConvertToProtoSong(song *Info) *proto.SongData {
	return &proto.SongData{
		Id:    song.id,
		Port:  int32(song.port),
		Songs: song.songs,
	}
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}
func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}

func (s *grpcFileServiceServer) saveFile(id, path string, buffer bytes.Buffer) error {
	fp := filepath.Join("assets", path, id)
	err := s.store.SaveChunk(fp, buffer)
	if err != nil {
		return err
	}
	return nil
}
