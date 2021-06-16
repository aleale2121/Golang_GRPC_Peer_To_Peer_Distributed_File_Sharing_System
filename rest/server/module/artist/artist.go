package module

import (
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/artist"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator/view"
	"net/http"
)

type UseCase interface {
	CreateArtist(artist *model.Artist) (*constant.SuccessData, *constant.ErrorData)
	Artists(sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData)
	Artist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData)
	GetArtistLikeCount(id string) (*constant.SuccessData, *constant.ErrorData)
	UpdateArtist(artist *model.Artist) (*constant.SuccessData, *constant.ErrorData)
	DeleteArtist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData)
	LikeArtist(albumLike *model.ArtistLikes) (*constant.SuccessData, *constant.ErrorData)
}
type Service struct {
	artistRepoService artist.ArtistsService
}

func (s Service) GetArtistLikeCount(id string) (*constant.SuccessData, *constant.ErrorData) {
	_, err := s.artistRepoService.Artist(uuid.FromStringOrNil(id))
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Artist With Given ID not found",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: struct {
			ArtistLikesCount int64 `json:"artist_likes_count"`
		}{
			ArtistLikesCount: s.artistRepoService.GetArtistLikeCount(uuid.FromStringOrNil(id)),
		},
	}, nil
}

func (s Service) LikeArtist(artistLike *model.ArtistLikes) (*constant.SuccessData, *constant.ErrorData) {
	_, err := s.artistRepoService.Artist(artistLike.ArtistId)

	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Artist With Given ID not found",
		}
	}
	isUnLiked, err := s.artistRepoService.LikeArtist(artistLike)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error  While Updating Like",
		}
	}
	if isUnLiked {
		return &constant.SuccessData{
			Code: http.StatusNoContent,
			Data: "Artist Like Undid",
		}, nil
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Artist Liked",
	}, nil
}

func (s Service) CreateArtist(artist *model.Artist) (*constant.SuccessData, *constant.ErrorData) {
	if artist.FirstName == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Artist Name Cannot Be Empty",
		}
	}
	_, err := s.artistRepoService.Artist(artist.ID)
	if err == nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Artist Already Created",
		}
	}
	_, err = s.artistRepoService.CreateArtist(artist)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Creating Artist",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Artist Created",
	}, nil
}

func (s Service) Artists(sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData) {
	isValid := IsValidFilter(sortKey) && (constant.IsValidSort(sort))

	if !isValid {
		sort = ""
		sortKey = ""
	}
	artists, p, err := s.artistRepoService.Artists(sort, sortKey, pageNo, maxPerPage)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Getting Artists",
		}
	}

	pageView := view.New(p)

	return &constant.SuccessData{
		Code: 200,
		Data: constant.PaginatedData{
			MetaData: constant.CreateMetaData(
				p.Page(), maxPerPage, p.PageNums(), p.Nums(),
				constant.LinksData{
					Self:     constant.GetFormattedLinkType1(isValid, "v1/artists", pageNo, maxPerPage, sort, sortKey),
					First:    constant.GetFormattedLinkType1(isValid, "v1/artists", 0, maxPerPage, sort, sortKey),
					Previous: constant.GetFormattedLinkType1(isValid, "v1/artists", pageView.Prev(), maxPerPage, sort, sortKey),
					Next:     constant.GetFormattedLinkType1(isValid, "v1/artists", pageView.Next(), maxPerPage, sort, sortKey),
					Last:     constant.GetFormattedLinkType1(isValid, "v1/artists", pageView.Last(), maxPerPage, sort, sortKey),
				}),
			Data: artists,
		},
	}, nil
}

func (s Service) Artist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData) {
	r, err := s.artistRepoService.Artist(id)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Artist With The Given Id not Found",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: r,
	}, nil
}

func (s Service) UpdateArtist(artist *model.Artist) (*constant.SuccessData, *constant.ErrorData) {
	if artist.FirstName == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Artist Name Cannot Be Empty",
		}
	}
	_, err := s.artistRepoService.Artist(artist.ID)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Artist With The Given ID not Found",
		}
	}
	_, err = s.artistRepoService.UpdateArtist(artist)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Updating Artist",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Artist Updated",
	}, nil
}

func (s Service) DeleteArtist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData) {
	_, err := s.artistRepoService.Artist(id)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Artist With The Given ID not Found",
		}
	}
	_, err = s.artistRepoService.DeleteArtist(id)
	if err != nil {

		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Deleting Artist",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Artist Deleted",
	}, nil
}

func NewArtistService(artistRepoService artist.ArtistsService) UseCase {
	return &Service{
		artistRepoService: artistRepoService,
	}
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
