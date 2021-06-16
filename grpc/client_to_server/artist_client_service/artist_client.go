package artist_client_service

import (
	"context"
	proto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/artist"
	"google.golang.org/grpc"
	"time"
)

type ArtistClient struct {
	service proto.ArtistServiceClient
}

func NewArtistClient(rc *grpc.ClientConn) *ArtistClient {
	return &ArtistClient{
		proto.NewArtistServiceClient(rc),
	}
}
func (artistClient *ArtistClient) CreateArtist(artist *proto.Artist) (string, error) {
	req := &proto.CreateArtistRequest{Artist: artist}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := artistClient.service.CreateArtist(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}
func (artistClient *ArtistClient) GetArtist(id string) (*proto.Artist, error) {
	req := &proto.GetArtistRequest{Id: id}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := artistClient.service.GetArtist(ctx, req)
	if err != nil {

		return nil, err
	}
	return res.Artist, nil

}
func (artistClient *ArtistClient) GetAllArtists(sort string, sortKey string, page int64, maxPerPage int64) (*proto.PaginatedArtistData, error) {

	res, err := artistClient.service.GetAllArtists(context.Background(),
		&proto.GetArtistsRequest{Sort: sort, SortKey: sortKey, Page: page, MaxPerPage: maxPerPage})
	if err != nil {

		return nil, err
	}

	return res.Artists, err
}
func (artistClient *ArtistClient) UpdateArtist(artist *proto.Artist) error {
	req := &proto.UpdateArtistRequest{Artist: artist}

	_, err := artistClient.service.UpdateArtist(context.Background(), req)
	if err != nil {

		return err
	}

	return nil
}
func (artistClient *ArtistClient) DeleteArtist(id string) error {
	req := &proto.DeleteArtistRequest{Id: id}
	_, err := artistClient.service.DeleteArtist(context.Background(), req)
	if err != nil {

		return err
	}

	return nil
}
func (artistClient *ArtistClient) LikeArtist(userId, artistId string) (bool, error) {
	req := &proto.LikeArtistRequest{
		UserId:   userId,
		ArtistId: artistId,
	}
	isLiked, err := artistClient.service.LikeArtist(context.Background(), req)
	if err != nil {
		return false, err
	}
	return isLiked.IsLiked, nil

}
func (artistClient *ArtistClient) GetArtistLikesCount(artistId string) (int64, error) {
	req := &proto.GetLikesCountRequest{
		ArtistId: artistId,
	}
	resp, err := artistClient.service.GetArtistLikesCount(context.Background(), req)
	if err != nil {

		return 0, err
	}

	return resp.Count, nil
}
