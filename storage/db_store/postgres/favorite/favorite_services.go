package favorite

import (
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator"
)

type FavoritesService interface {
	CreateFavorite(favorite *model2.Favorite) (*model2.Favorite, error)
	GetFavoriteById(favoriteId uuid.UUID) (*model2.Favorite, error)
	GetUserFavoriteSongs(userId uuid.UUID, sort string, sortKey string, pageNo int, maxPerPage int) ([]model2.Favorite, *paginator.Paginator, error)
	GetUserFavoriteSong(userId, songId uuid.UUID) (*model2.Favorite, error)
	RemoveSongFromUserFavorite(userId, songId uuid.UUID) error
	UpdateFavoriteName(favorite model2.Favorite) (*model2.Favorite, error)
}
