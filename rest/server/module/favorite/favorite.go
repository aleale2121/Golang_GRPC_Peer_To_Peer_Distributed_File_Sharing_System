package module

import (
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/favorite"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator/view"
	"net/http"
)

type UseCase interface {
	CreateFavorite(favorite *model2.Favorite) (*constant.SuccessData, *constant.ErrorData)
	UpdateFavorite(title, favId string) (*constant.SuccessData, *constant.ErrorData)
	GetUserFavoriteSongs(userId uuid.UUID, sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData)
	RemoveSongFromUserFavorite(userId, songId uuid.UUID) (*constant.SuccessData, *constant.ErrorData)
}
type Service struct {
	favoriteRepoService favorite.FavoritesService
}

func NewFavoriteService(favoriteRepoService favorite.FavoritesService) UseCase {
	return &Service{
		favoriteRepoService: favoriteRepoService,
	}
}

func (s Service) CreateFavorite(favorite *model2.Favorite) (*constant.SuccessData, *constant.ErrorData) {
	if favorite.Title == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Favorite title cannot be empty",
		}
	}
	_, err := s.favoriteRepoService.GetUserFavoriteSong(favorite.UserId, favorite.SongId)
	if err == nil {
		err = s.favoriteRepoService.RemoveSongFromUserFavorite(favorite.UserId, favorite.SongId)
		if err != nil {
			return nil, &constant.ErrorData{
				Code:  http.StatusInternalServerError,
				Title: "Error Occurred While removing song from favorites",
			}
		}
		return &constant.SuccessData{
			Code: http.StatusNoContent,
			Data: "Song Removed from favorites",
		}, nil
	}

	_, err = s.favoriteRepoService.CreateFavorite(favorite)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While removing song from favorites",
		}
	}

	return &constant.SuccessData{
		Code: 200,
		Data: "Song Added To favorite list",
	}, nil
}

func (s Service) UpdateFavorite(title, favID string) (*constant.SuccessData, *constant.ErrorData) {
	if title == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Favorite title cannot be empty",
		}
	}
	fmt.Println(favID)
	_, err := s.favoriteRepoService.GetFavoriteById(uuid.FromStringOrNil(favID))
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "No Favorite Song Found with the given favorite id",
		}
	}

	_, err = s.favoriteRepoService.UpdateFavoriteName(model2.Favorite{
		ID:    uuid.FromStringOrNil(favID),
		Title: title,
	})
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While updating favorites",
		}
	}

	return &constant.SuccessData{
		Code: 200,
		Data: "Favorite Updated",
	}, nil
}

func (s Service) GetUserFavoriteSongs(userId uuid.UUID, sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData) {
	isValid := IsValidFilter(sortKey) && (constant.IsValidSort(sort))

	if !isValid {
		sort = ""
		sortKey = ""
	}
	albums, p, err := s.favoriteRepoService.GetUserFavoriteSongs(userId, sort, sortKey, pageNo, maxPerPage)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Getting Favorite Songs",
		}
	}
	pageView := view.New(p)

	return &constant.SuccessData{
		Code: 200,
		Data: constant.PaginatedData{
			MetaData: constant.CreateMetaData(
				p.Page(), maxPerPage, p.PageNums(), p.Nums(),
				constant.LinksData{
					Self:     constant.GetFormattedLinkType1(isValid, "v1/favorites", pageNo, maxPerPage, sort, sortKey),
					First:    constant.GetFormattedLinkType1(isValid, "v1/favorites", 0, maxPerPage, sort, sortKey),
					Previous: constant.GetFormattedLinkType1(isValid, "v1/favorites", pageView.Prev(), maxPerPage, sort, sortKey),
					Next:     constant.GetFormattedLinkType1(isValid, "v1/favorites", pageView.Next(), maxPerPage, sort, sortKey),
					Last:     constant.GetFormattedLinkType1(isValid, "v1/favorites", pageView.Last(), maxPerPage, sort, sortKey),
				}),
			Data: albums,
		},
	}, nil
}

func (s Service) RemoveSongFromUserFavorite(userId, songId uuid.UUID) (*constant.SuccessData, *constant.ErrorData) {
	_, err := s.favoriteRepoService.GetUserFavoriteSong(userId, songId)
	if err == nil {
		err = s.favoriteRepoService.RemoveSongFromUserFavorite(userId, songId)
		if err != nil {
			return nil, &constant.ErrorData{
				Code:  http.StatusInternalServerError,
				Title: "Error Occurred While removing song from favorite",
			}
		}
		return &constant.SuccessData{
			Code: http.StatusNoContent,
			Data: "Song Removed from favorite list",
		}, nil
	}

	return nil, &constant.ErrorData{
		Code:  http.StatusNotFound,
		Title: "No favorite song found for the given user",
	}

}
func IsValidFilter(key string) bool {
	keys := [1]string{"title"}
	for i := 0; i < len(keys); i++ {
		if keys[i] == key {
			return true
		}
	}
	return false
}
