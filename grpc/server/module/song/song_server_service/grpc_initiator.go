package song_server_service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	proto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/song"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/artist"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/song"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/file_store"
	"github.com/gabriel-vasile/mimetype"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator/view"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type grpcSongServiceServer struct {
	songRepoService   song.SongsService
	artistRepoService artist.ArtistsService
	store             file_store.Storage
	proto.UnimplementedSongServiceServer
}


func NewGrpcSongServer(songRepoService song.SongsService,
	artistRepoService artist.ArtistsService,
	store              file_store.Storage, ) proto.SongServiceServer {
	return &grpcSongServiceServer{
		songRepoService:   songRepoService,
		artistRepoService: artistRepoService,
		store: store,
	}
}
func (s grpcSongServiceServer) CreateSong(stream proto.SongService_CreateSongServer) error {

	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive image info")
	}

	songInfo := req.GetSongInfo()

	if songInfo.Title == "" {
		return  status.Errorf(codes.InvalidArgument, "Song Title Cannot Be Empty")
	}
	if songInfo.ArtistId == "" {

		return  status.Errorf(codes.InvalidArgument, "Artist Id Cannot Be Empty")
	}
	_, err = s.artistRepoService.Artist(uuid.FromStringOrNil(songInfo.ArtistId))
	if err != nil {

		return status.Errorf(codes.NotFound, "Artist With Given User ID not found")
	}
	coverImageData := bytes.Buffer{}

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

		coverImageChunkData := req.GetSongCover_ImageChunkData()

		_, err = coverImageData.Write(coverImageChunkData)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write cover image chunk data: %v", err)
		}

	}
	mime := mimetype.Detect(coverImageData.Bytes())
	if !mimetype.EqualsAny(mime.String(),"image/jpeg","image/pjpeg",
		"image/png", "image/tiff","image/x-tiff","image/vnd.wap.wbmp"){
		return status.Errorf(codes.InvalidArgument, "the cover image you upload is not image")
	}
	id := uuid.NewV4().String()
	coverImageId := uuid.NewV4().String()+mime.Extension()
	err = s.saveFile(id,coverImageId,coverImageData)
	if err != nil {
		return  fmt.Errorf("cannot write song cover to file: %w", err)
	}

	sg, err := s.songRepoService.CreateSong(&model2.Song{
		ArtistId:      uuid.FromStringOrNil(songInfo.ArtistId),
		Title:         songInfo.Title,
		CoverImageUrl:  id + "/" + coverImageId,
		SongUrl:         "",
		Views:         0,
		Duration:      int(songInfo.Duration),
		ReleaseAt:     time.Now(),
	})
	if err != nil {
		return  status.Errorf(codes.Internal, "Error Occurred While Creating Song")
	}
	res := &proto.CreateSongResponse{
		Id:   sg.ID.String(),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}
	log.Println("song saved")
	return nil
}


func (s grpcSongServiceServer) UpdateSongUrl(stream proto.SongService_UpdateSongUrlServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive image info")
	}

	songId := req.GetSongId()
	songX, err := s.songRepoService.Song(uuid.FromStringOrNil(songId))
	if err != nil {
		return status.Errorf(codes.NotFound, "Artist With Given User ID not found")
	}
	songData := bytes.Buffer{}


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
		songChunkData := req.GetSongChunkData()

		_, err = songData.Write(songChunkData)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write song data: %v", err)
		}
	}
	mime := mimetype.Detect(songData.Bytes())
	if !mimetype.EqualsAny(mime.String(),"audio/basic","audio/L24",
		"audio/mid", "audio/mpeg","audio/mp4","audio/x-aiff","audio/x-mpegurl",
		"audio/ogg","audio/vorbis","audio/vnd.wav","text/plain"){
		return status.Errorf(codes.InvalidArgument, "the song you upload is not audio")
	}
	id:=strings.Split(songX.CoverImageUrl,"/")[0]
	songUrlId := uuid.NewV4().String()+mime.Extension()

	err = s.saveFile(id,songUrlId,songData)
	if err != nil {
		return  fmt.Errorf("cannot write ssong to file: %w", err)
	}
	songX.SongUrl=id+"/"+songUrlId
	_, err = s.songRepoService.UpdateSong(songX)
	if err != nil {
		return  status.Errorf(codes.Internal, "Error Occurred While Updating Song Url")
	}
	err = stream.SendAndClose(&proto.UpdateSongUrlResponse{})
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}
	log.Println("song updated")
	return nil
}


func (s grpcSongServiceServer) GetArtistSongs(_ context.Context, request *proto.GetArtistSongsRequest) (*proto.GetSongsResponse, error) {
	isValid := IsValidFilter(request.SortKey) && (constant.IsValidSort(request.Sort))
	if !isValid {
		request.Sort = ""
		request.SortKey = ""
	}
	albums, p, err := s.songRepoService.SongByArtistID(request.ArtistId, request.Sort, request.SortKey, int(request.Page), int(request.MaxPerPage))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Error Occurred While Getting Songs")
	}
	pageView := view.New(p)
	return &proto.GetSongsResponse{
		Songs: ConvertToProtoPaginatedSongsData(
			albums,
			CreateMetaData(
				p.Page(), int(request.MaxPerPage), p.PageNums(), p.Nums(),
				constant.LinksData{
					Self:     constant.GetFormattedLinkType1(isValid, "songs", int(request.Page), int(request.MaxPerPage), request.Sort, request.SortKey),
					First:    constant.GetFormattedLinkType1(isValid, "songs", 0, int(request.MaxPerPage), request.Sort, request.SortKey),
					Previous: constant.GetFormattedLinkType1(isValid, "songs", pageView.Prev(), int(request.MaxPerPage), request.Sort, request.SortKey),
					Next:     constant.GetFormattedLinkType1(isValid, "songs", pageView.Next(), int(request.MaxPerPage), request.Sort, request.SortKey),
					Last:     constant.GetFormattedLinkType1(isValid, "songs", pageView.Last(), int(request.MaxPerPage), request.Sort, request.SortKey),
				},
			),
		)}, nil
}

func (s grpcSongServiceServer) GetSongViewCount(_ context.Context, request *proto.GetViewsCountRequest) (*proto.GetViewsCountResponse, error) {
	_, err := s.songRepoService.Song(uuid.FromStringOrNil(request.SongId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Song With The Given ID not Found")

	}
	views, err := s.songRepoService.GetSongViewsCount(uuid.FromStringOrNil(request.SongId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error Occurred While getting song view count")
	}
	return &proto.GetViewsCountResponse{
		Count: views,
	}, nil
}

func (s grpcSongServiceServer) IncreaseSongViewCount(_ context.Context, request *proto.IncreaseSongViewRequest) (*proto.IncreaseSongViewResponse, error) {
	_, err := s.songRepoService.Song(uuid.FromStringOrNil(request.SongId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Song With The Given ID not Found")

	}
	err = s.songRepoService.IncreaseSongViews(request.SongId)
	if err != nil {

		return nil, status.Errorf(codes.Internal, "Error Occurred While updating view")

	}
	return &proto.IncreaseSongViewResponse{}, nil
}

func (s grpcSongServiceServer) LikeSong(_ context.Context, request *proto.LikeSongRequest) (*proto.LikeSongResponse, error) {
	_, err := s.songRepoService.Song(uuid.FromStringOrNil(request.SongId))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Song With The Given ID not Found")

	}

	isUnLiked, err := s.songRepoService.LikeSong(&model2.SongLikes{
		SongId: uuid.FromStringOrNil(request.SongId),
		UserId: uuid.FromStringOrNil(request.UserId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error Occurred While Updating Like")

	}
	if isUnLiked {
		return &proto.LikeSongResponse{IsLiked: true}, nil
	}
	return &proto.LikeSongResponse{IsLiked: false}, nil
}

func (s grpcSongServiceServer) GetSongLikesCount(_ context.Context, request *proto.GetLikesCountRequest) (*proto.GetLikesCountResponse, error) {
	_, err := s.songRepoService.Song(uuid.FromStringOrNil(request.SongId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Song With The Given Id not Found")
	}
	return &proto.GetLikesCountResponse{
		Count: s.songRepoService.GetSongLikeCount(uuid.FromStringOrNil(request.SongId))}, nil
}

func (s grpcSongServiceServer) GetSong(_ context.Context, request *proto.GetSongRequest) (*proto.GetSongResponse, error) {
	r, err := s.songRepoService.Song(uuid.FromStringOrNil(request.SongId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Song With The Given Id not Found")
	}
	return &proto.GetSongResponse{Song: ConvertToProtoSong(r)}, nil
}

func (s grpcSongServiceServer) GetAllSongs(_ context.Context, request *proto.GetSongsRequest) (*proto.GetSongsResponse, error) {

	isValid := IsValidFilter(request.SortKey) && (constant.IsValidSort(request.Sort))
	if !isValid {
		request.Sort = ""
		request.SortKey = ""
	}
	songs, p, err := s.songRepoService.Songs(request.Sort, request.SortKey, int(request.Page), int(request.MaxPerPage))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Error Occurred While Getting Songs")
	}

	pageView := view.New(p)
	return &proto.GetSongsResponse{
		Songs: ConvertToProtoPaginatedSongsData(
			songs,
			CreateMetaData(
				p.Page(), int(request.MaxPerPage), p.PageNums(), p.Nums(),
				constant.LinksData{
					Self:     constant.GetFormattedLinkType1(isValid, "songs", int(request.Page), int(request.MaxPerPage), request.Sort, request.SortKey),
					First:    constant.GetFormattedLinkType1(isValid, "songs", 0, int(request.MaxPerPage), request.Sort, request.SortKey),
					Previous: constant.GetFormattedLinkType1(isValid, "songs", pageView.Prev(), int(request.MaxPerPage), request.Sort, request.SortKey),
					Next:     constant.GetFormattedLinkType1(isValid, "songs", pageView.Next(), int(request.MaxPerPage), request.Sort, request.SortKey),
					Last:     constant.GetFormattedLinkType1(isValid, "songs", pageView.Last(), int(request.MaxPerPage), request.Sort, request.SortKey),
				},
			),
		)}, nil
}

func (s grpcSongServiceServer) UpdateSong(_ context.Context, request *proto.UpdateSongRequest) (*proto.UpdateSongResponse, error) {
	if request.Song.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Song Title Cannot Be Empty")
	}
	if request.Song.ArtistId == "" {

		return nil, status.Errorf(codes.InvalidArgument, "Artist Id Cannot Be Empty")
	}
	_, err := s.songRepoService.Song(uuid.FromStringOrNil(request.Song.ID))
	if err != nil {

		return nil, status.Errorf(codes.NotFound, "Song With Given User ID not found")
	}
	_, err = s.artistRepoService.Artist(uuid.FromStringOrNil(request.Song.ArtistId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Artist With Given User ID not found")
	}
	_, err = s.songRepoService.UpdateSong(ConvertProtoToSong(request.Song))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error Occurred While Creating Song")
	}
	return &proto.UpdateSongResponse{}, nil
}

func (s grpcSongServiceServer) DeleteSong(_ context.Context, request *proto.DeleteSongRequest) (*proto.DeleteSongResponse, error) {
	_, err := s.songRepoService.Song(uuid.FromStringOrNil(request.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Song With The Given ID not Found")

	}
	_, err = s.songRepoService.DeleteSong(uuid.FromStringOrNil(request.Id))
	if err != nil {

		return nil, status.Errorf(codes.Internal, "Error Occurred While Deleting Song")

	}
	return &proto.DeleteSongResponse{}, nil
}

func (s grpcSongServiceServer) DownloadFile(request *proto.DownloadFileRequest, server proto.SongService_DownloadFileServer) error {
	wd,_:=os.Getwd()
	fp := filepath.Join(wd,"assets", request.FileId)

	file, err := os.Open(fp)
	if err != nil {
		log.Fatal("cannot open coverImage: ", err)
	}
	defer file.Close()

	buff := make([]byte, 1024)
	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				break
			}
		}
		resp := &proto.DownloadFileResponse{
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

func ConvertProtoToArtist(artist *proto.Artist) *model2.Artist {
	return &model2.Artist{
		ID:        uuid.FromStringOrNil(artist.ID),
		FirstName: artist.FirstName,
		LastName:  artist.LastName,
		CreatedAt: artist.CreatedAt.AsTime(),
		UpdatedAt: artist.UpdatedAt.AsTime(),
	}
}

func ConvertToProtoArtist(artist *model2.Artist) *proto.Artist {
	return &proto.Artist{
		ID:        artist.ID.String(),
		FirstName: artist.FirstName,
		LastName:  artist.LastName,
		CreatedAt: timestamppb.New(artist.CreatedAt),
		UpdatedAt: timestamppb.New(artist.UpdatedAt),
	}
}

func ConvertToProtoPaginatedSongsData(songs []model2.Song, metaData *proto.SongMetaData) *proto.PaginatedSongData {
	return &proto.PaginatedSongData{
		MetaData: metaData,
		Data:     ConvertToArrayOfProtoSongs(songs),
	}
}

func CreateMetaData(page int, perPage int, pageCount int, totalCount int, data constant.LinksData) *proto.SongMetaData {
	return &proto.SongMetaData{
		Page:       int64(page),
		PerPage:    int64(perPage),
		PageCount:  int64(pageCount),
		TotalCount: int64(totalCount),
		Links: []*proto.SongLink{
			{
				Key:   "self",
				Value: data.Self,
			},
			{
				Key:   "previous",
				Value: data.Self,
			},
			{
				Key:   "next",
				Value: data.Self,
			},
			{
				Key:   "self",
				Value: data.Self,
			},
			{
				Key:   "last",
				Value: data.Self,
			},
		},
	}
}

func ConvertProtoToSong(song *proto.Song) *model2.Song {
	return &model2.Song{
		ID:            uuid.FromStringOrNil(song.ID),
		ArtistId:      uuid.FromStringOrNil(song.ArtistId),
		Artist:        *ConvertProtoToArtist(song.Artist),
		Title:         song.Title,
		CoverImageUrl: song.CoverImageUrl,
		Views:         int(song.Views),
		Duration:      int(song.Duration),
		ReleaseAt:     song.ReleasedAt.AsTime(),
		CreatedAt:     song.CreatedAt.AsTime(),
		UpdatedAt:     song.UpdatedAt.AsTime(),
	}
}

func ConvertToProtoSong(song *model2.Song) *proto.Song {
	return &proto.Song{
		ID:            song.ID.String(),
		ArtistId:      song.ArtistId.String(),
		Artist:        ConvertToProtoArtist(&song.Artist),
		Title:         song.Title,
		CoverImageUrl: song.CoverImageUrl,
		SongUrl: song.SongUrl,
		Views:         int64(song.Views),
		Duration:      int64(song.Duration),
		ReleasedAt:    timestamppb.New(song.ReleaseAt),
		CreatedAt:     timestamppb.New(song.CreatedAt),
		UpdatedAt:     timestamppb.New(song.UpdatedAt),
	}
}

func ConvertToArrayOfProtoSongs(songs []model2.Song) []*proto.Song {
	var protoSongs []*proto.Song
	for i := 0; i < len(songs); i++ {
		protoSongs = append(protoSongs, ConvertToProtoSong(&songs[i]))
	}
	return protoSongs
}

func IsValidFilter(key string) bool {
	keys := [4]string{"title", "views", "duration", "release_at"}
	for i := 0; i < len(keys); i++ {
		if keys[i] == key {
			return true
		}
	}
	return false
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

func (s *grpcSongServiceServer) saveFile(id, path string, buffer bytes.Buffer) error {

	fp := filepath.Join("assets", id, path)
	err := s.store.SaveChunk(fp, buffer)
	if err != nil {
		return err
	}
	return nil
}
