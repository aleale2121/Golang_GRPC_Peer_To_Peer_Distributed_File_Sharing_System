package artist

import (
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator"
)

type ArtistsService interface {
	CreateArtist(artist *model2.Artist) (*model2.Artist, error)
	Artists(sort string, sortKey string, pageNo int, maxPerPage int) ([]model2.Artist, *paginator.Paginator, error)
	Artist(id uuid.UUID) (*model2.Artist, error)
	UpdateArtist(artist *model2.Artist) (*model2.Artist, error)
	DeleteArtist(id uuid.UUID) (*model2.Artist, error)
	LikeArtist(artistLike *model2.ArtistLikes) (bool, error)
	GetArtistLikeCount(id uuid.UUID) int64
}
