package song_client_service

import (
	"bufio"
	"bytes"
	"context"
	proto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/song"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/file_store"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type SongClient struct {
	service proto.SongServiceClient
	store   file_store.Storage

}

func NewSongClient(rc *grpc.ClientConn,	store file_store.Storage) *SongClient {
	return &SongClient{
		service: proto.NewSongServiceClient(rc),
		store: store,
	}
}
func (songClient *SongClient) CreateSong(song *proto.SongCreateInfo,coverImagePath,songPath string) (string, error) {
	coverImage, err := os.Open(coverImagePath)
	if err != nil {
		log.Fatal("cannot open coverImage: ", err)
	}
	defer coverImage.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	stream, err := songClient.service.CreateSong(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req := &proto.CreateSongRequest{
		Data: &proto.CreateSongRequest_SongInfo{
			SongInfo: song},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send song info to server: ", err, stream.RecvMsg(nil))
	}

	readerCoverImage := bufio.NewReader(coverImage)
	coverImageBuffers := make([]byte, 1024)

	for {
		n, err := readerCoverImage.Read(coverImageBuffers)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to coverImageBuffers: ", err)
		}

		req := &proto.CreateSongRequest{
			Data: &proto.CreateSongRequest_SongCover_ImageChunkData{
				SongCover_ImageChunkData: coverImageBuffers[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send cover chunk to server: ", err, stream.RecvMsg(nil))
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	stream2, err := songClient.service.UpdateSongUrl(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}
	songAudio, err := os.Open(songPath)
	if err != nil {
		log.Fatal("cannot open song: ", err)
	}
	defer songAudio.Close()

	req2 := &proto.UpdateSongUrlRequest{
		Data: &proto.UpdateSongUrlRequest_SongId{SongId: res.Id},
	}

	err = stream2.Send(req2)
	if err != nil {
		log.Fatal("cannot send update song url to server: ", err,
			stream2.RecvMsg(nil))
	}
	readerSong := bufio.NewReader(songAudio)
	songBuffers := make([]byte, 1024)

	for {
		n, err := readerSong.Read(songBuffers)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to song buffers: ", err)
		}

		req2 := &proto.UpdateSongUrlRequest{
			Data: &proto.UpdateSongUrlRequest_SongChunkData{
				SongChunkData: songBuffers[:n],
			},
		}

		err = stream2.Send(req2)
		if err != nil {
			log.Fatal("cannot send song chunk to server: ", err, stream2.RecvMsg(nil))
		}
	}
	_, err = stream2.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	return res.Id, err
}
func (songClient *SongClient) GetSong(id string) (*proto.Song, error) {
	req := &proto.GetSongRequest{SongId: id}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := songClient.service.GetSong(ctx, req)
	if err != nil {

		return nil, err
	}
	return res.Song, nil

}
func (songClient *SongClient) GetAllSongs(sort string, sortKey string, page int64, maxPerPage int64) (*proto.PaginatedSongData, error) {

	res, err := songClient.service.GetAllSongs(context.Background(),
		&proto.GetSongsRequest{SortKey: sortKey, Sort: sort, Page: page, MaxPerPage: maxPerPage})
	if err != nil {

		return nil, err
	}

	return res.Songs, err
}
func (songClient *SongClient) GetArtistSongs(artistID string, sort string, sortKey string, page int64, maxPerPage int64) (*proto.PaginatedSongData, error) {

	res, err := songClient.service.GetArtistSongs(context.Background(),
		&proto.GetArtistSongsRequest{ArtistId: artistID, SortKey: sortKey, Sort: sort, Page: page, MaxPerPage: maxPerPage})
	if err != nil {

		return nil, err
	}

	return res.Songs, err
}
func (songClient *SongClient) UpdateSong(song *proto.Song) error {
	req := &proto.UpdateSongRequest{Song: song}

	_, err := songClient.service.UpdateSong(context.Background(), req)
	if err != nil {

		return err
	}

	return nil
}
func (songClient *SongClient) DeleteSong(id string) error {
	req := &proto.DeleteSongRequest{Id: id}
	_, err := songClient.service.DeleteSong(context.Background(), req)
	if err != nil {

		return err
	}

	return nil
}
func (songClient *SongClient) LikeSong(userId, songId string) (bool, error) {
	req := &proto.LikeSongRequest{
		UserId: userId,
		SongId: songId,
	}
	isLiked, err := songClient.service.LikeSong(context.Background(), req)
	if err != nil {
		return false, err
	}
	return isLiked.IsLiked, nil

}
func (songClient *SongClient) GetSongLikesCount(songId string) (int64, error) {
	req := &proto.GetLikesCountRequest{
		SongId: songId,
	}
	resp, err := songClient.service.GetSongLikesCount(context.Background(), req)
	if err != nil {

		return 0, err
	}

	return resp.Count, nil
}

func (songClient *SongClient) GetSongViewCount(songId string) (int64, error) {
	req := &proto.GetViewsCountRequest{
		SongId: songId,
	}
	resp, err := songClient.service.GetSongViewCount(context.Background(), req)
	if err != nil {

		return 0, err
	}

	return resp.Count, nil
}
func (songClient *SongClient) IncreaseAlbumViewCount(songId string) error {
	req := &proto.IncreaseSongViewRequest{
		SongId: songId,
	}
	_, err := songClient.service.IncreaseSongViewCount(context.Background(), req)
	if err != nil {

		return err
	}

	return nil
}
func (songClient *SongClient) DownloadFile(fileId string) error {
	_ = &proto.DownloadFileRequest{
		FileId: fileId,
	}

	return nil
}

func (songClient *SongClient) saveFile(id, path string, buffer bytes.Buffer) error {

	fp := filepath.Join("downloads", id, path)
	err := songClient.store.SaveChunk(fp, buffer)
	if err != nil {
		return err
	}
	return nil
}
