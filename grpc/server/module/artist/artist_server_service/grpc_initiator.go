package artist_server_service

import (
	"context"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	proto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/artist"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/artist"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator/view"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcArtistServiceServer struct {
	artistRepoService artist.ArtistsService
}

func (s grpcArtistServiceServer) CreateArtist(_ context.Context, request *proto.CreateArtistRequest) (*proto.CreateArtistResponse, error) {

	if request.Artist.FirstName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Artist Name Cannot Be Empty")
	}
	_, err := s.artistRepoService.Artist(uuid.FromStringOrNil(request.Artist.ID))
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "Artist Already Created")
	}
	art, err := s.artistRepoService.CreateArtist(ConvertProtoToArtist(request.Artist))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error Occurred While Creating Artist")
	}
	return &proto.CreateArtistResponse{Id: art.ID.String()}, nil
}

func (s grpcArtistServiceServer) GetArtist(_ context.Context, request *proto.GetArtistRequest) (*proto.GetArtistResponse, error) {
	r, err := s.artistRepoService.Artist(uuid.FromStringOrNil(request.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Artist With The Given Id not Found")
	}
	return &proto.GetArtistResponse{Artist: ConvertToProtoArtist(r)}, nil
}

func (s grpcArtistServiceServer) GetAllArtists(_ context.Context, request *proto.GetArtistsRequest) (*proto.GetArtistsResponse, error) {
	isValid := IsValidFilter(request.SortKey) && (constant.IsValidSort(request.Sort))
	if !isValid {
		request.Sort = ""
		request.SortKey = ""
	}
	artists, p, err := s.artistRepoService.Artists(request.Sort, request.SortKey, int(request.Page), int(request.MaxPerPage))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Error Occurred While Getting Artists")
	}
	pageView := view.New(p)
	return &proto.GetArtistsResponse{
		Artists: ConvertToProtoPaginatedArtistsData(
			artists,
			CreateMetaData(
				p.Page(), int(request.MaxPerPage), p.PageNums(), p.Nums(),
				constant.LinksData{
					Self:     constant.GetFormattedLinkType1(isValid, "artist", int(request.Page), int(request.MaxPerPage), request.Sort, request.SortKey),
					First:    constant.GetFormattedLinkType1(isValid, "artist", 0, int(request.MaxPerPage), request.Sort, request.SortKey),
					Previous: constant.GetFormattedLinkType1(isValid, "artist", pageView.Prev(), int(request.MaxPerPage), request.Sort, request.SortKey),
					Next:     constant.GetFormattedLinkType1(isValid, "artist", pageView.Next(), int(request.MaxPerPage), request.Sort, request.SortKey),
					Last:     constant.GetFormattedLinkType1(isValid, "artist", pageView.Last(), int(request.MaxPerPage), request.Sort, request.SortKey),
				},
			),
		)}, nil
}

func (s grpcArtistServiceServer) UpdateArtist(_ context.Context, request *proto.UpdateArtistRequest) (*proto.UpdateArtistResponse, error) {
	if request.Artist.FirstName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Artist Name Cannot Be Empty")
	}
	_, err := s.artistRepoService.Artist(uuid.FromStringOrNil(request.Artist.ID))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Artist With The Given ID not Found")
	}
	_, err = s.artistRepoService.UpdateArtist(ConvertProtoToArtist(request.Artist))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error Occurred While Updating Artist")
	}
	return &proto.UpdateArtistResponse{}, nil
}

func (s grpcArtistServiceServer) DeleteArtist(_ context.Context, request *proto.DeleteArtistRequest) (*proto.DeleteArtistResponse, error) {
	_, err := s.artistRepoService.Artist(uuid.FromStringOrNil(request.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Artist With The Given ID not Found")

	}
	_, err = s.artistRepoService.DeleteArtist(uuid.FromStringOrNil(request.Id))
	if err != nil {

		return nil, status.Errorf(codes.Internal, "Error Occurred While Deleting Artist")

	}
	return &proto.DeleteArtistResponse{}, nil
}

func (s grpcArtistServiceServer) LikeArtist(_ context.Context, request *proto.LikeArtistRequest) (*proto.LikeArtistResponse, error) {
	_, err := s.artistRepoService.Artist(uuid.FromStringOrNil(request.ArtistId))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Artist With The Given ID not Found")

	}

	isUnLiked, err := s.artistRepoService.LikeArtist(&model2.ArtistLikes{
		ArtistId: uuid.FromStringOrNil(request.ArtistId),
		UserId:   uuid.FromStringOrNil(request.UserId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error Occurred While Updating Like")

	}
	if isUnLiked {
		return &proto.LikeArtistResponse{IsLiked: true}, nil
	}
	return &proto.LikeArtistResponse{IsLiked: false}, nil
}

func (s grpcArtistServiceServer) GetArtistLikesCount(_ context.Context, request *proto.GetLikesCountRequest) (*proto.GetLikesCountResponse, error) {
	_, err := s.artistRepoService.Artist(uuid.FromStringOrNil(request.ArtistId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Artist With The Given ID not Found")

	}
	return &proto.GetLikesCountResponse{
		Count: s.artistRepoService.GetArtistLikeCount(uuid.FromStringOrNil(request.ArtistId)),
	}, nil
}

func NewGrpcArtistAServer(artistRepoService artist.ArtistsService) proto.ArtistServiceServer {
	return &grpcArtistServiceServer{artistRepoService: artistRepoService}
}

func ConvertProtoToArtist(artist *proto.Artist) *model2.Artist {
	return &model2.Artist{
		ID:        uuid.FromStringOrNil(artist.ID),
		FirstName: artist.FirstName,
		LastName:  artist.LastName,
		Image:     artist.Image,
		Email:     artist.Email,
		CreatedAt: artist.CreatedAt.AsTime(),
		UpdatedAt: artist.UpdatedAt.AsTime(),
	}
}
func ConvertToProtoArtist(artist *model2.Artist) *proto.Artist {
	return &proto.Artist{
		ID:        artist.ID.String(),
		FirstName: artist.FirstName,
		LastName:  artist.LastName,
		Image:     artist.Image,
		Email:     artist.Email,
		CreatedAt: timestamppb.New(artist.CreatedAt),
		UpdatedAt: timestamppb.New(artist.UpdatedAt),
	}
}
func ConvertToProtoPaginatedArtistsData(artist []model2.Artist, metaData *proto.MetaData) *proto.PaginatedArtistData {
	return &proto.PaginatedArtistData{
		MetaData: metaData,
		Data:     ConvertToArrayOfProtoArtists(artist),
	}
}
func CreateMetaData(page int, perPage int, pageCount int, totalCount int, data constant.LinksData) *proto.MetaData {
	return &proto.MetaData{
		Page:       int64(page),
		PerPage:    int64(perPage),
		PageCount:  int64(pageCount),
		TotalCount: int64(totalCount),
		Links: []*proto.Link{
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
func ConvertToArrayOfProtoArtists(artists []model2.Artist) []*proto.Artist {
	var protoArtists []*proto.Artist
	for i := 0; i < len(artists); i++ {
		protoArtists = append(protoArtists, ConvertToProtoArtist(&artists[i]))
	}
	return protoArtists
}
func IsValidFilter(key string) bool {
	keys := [2]string{"first_name", "last_name"}
	for i := 0; i < len(keys); i++ {
		if keys[i] == key {
			return true
		}
	}
	return false
}
