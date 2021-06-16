package module

import (
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/artist"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/song"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator/view"
	"net/http"
)

type UseCase interface {
	CreateSong(album *model2.Song) (*constant.SuccessData, *constant.ErrorData)
	Songs(sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData)
	ArtistSongs(artistID string, sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData)
	Song(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData)
	GetSongLikeCount(id string) (*constant.SuccessData, *constant.ErrorData)
	UpdateSong(album *model2.Song) (*constant.SuccessData, *constant.ErrorData)
	DeleteSong(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData)
	LikeSong(songLike *model2.SongLikes) (*constant.SuccessData, *constant.ErrorData)
	IncreaseSongViews(id string) (*constant.SuccessData, *constant.ErrorData)
	GetSongViewsCount(id string) (*constant.SuccessData, *constant.ErrorData)
}
type Service struct {
	songRepoService   song.SongsService
	artistRepoService artist.ArtistsService
}

func (s Service) IncreaseSongViews(id string) (*constant.SuccessData, *constant.ErrorData) {
	_, err := s.songRepoService.Song(uuid.FromStringOrNil(id))
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Song With Given ID not found",
		}
	}
	err = s.songRepoService.IncreaseSongViews(id)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While updating view",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Song View Successfully Updated",
	}, nil
}

func (s Service) GetSongViewsCount(id string) (*constant.SuccessData, *constant.ErrorData) {
	_, err := s.songRepoService.Song(uuid.FromStringOrNil(id))
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Song With Given ID not found",
		}
	}
	views, err := s.songRepoService.GetSongViewsCount(uuid.FromStringOrNil(id))
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While getting song view count",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: struct {
			SongViewsCount int64 `json:"song_views_count"`
		}{
			SongViewsCount: views,
		},
	}, nil
}

func (s Service) ArtistSongs(artistID string, sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData) {
	isValid := IsValidFilter(sortKey) && (constant.IsValidSort(sort))

	if !isValid {
		sort = ""
		sortKey = ""
	}
	songs, p, err := s.songRepoService.SongByArtistID(artistID, sort, sortKey, pageNo, maxPerPage)
	artistSongLink := fmt.Sprintf("v1/artist/songs?artist_id=%s", artistID)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Getting Songs",
		}
	}
	pageView := view.New(p)

	return &constant.SuccessData{
		Code: 200,
		Data: constant.PaginatedData{
			MetaData: constant.CreateMetaData(
				p.Page(), maxPerPage, p.PageNums(), p.Nums(),
				constant.LinksData{
					Self:     constant.GetFormattedLinkType2(isValid, artistSongLink, pageNo, maxPerPage, sort, sortKey),
					First:    constant.GetFormattedLinkType2(isValid, artistSongLink, 0, maxPerPage, sort, sortKey),
					Previous: constant.GetFormattedLinkType2(isValid, artistSongLink, pageView.Prev(), maxPerPage, sort, sortKey),
					Next:     constant.GetFormattedLinkType2(isValid, artistSongLink, pageView.Next(), maxPerPage, sort, sortKey),
					Last:     constant.GetFormattedLinkType2(isValid, artistSongLink, pageView.Last(), maxPerPage, sort, sortKey),
				}),
			Data: songs,
		},
	}, nil
}

func (s Service) GetSongLikeCount(id string) (*constant.SuccessData, *constant.ErrorData) {
	_, err := s.songRepoService.Song(uuid.FromStringOrNil(id))
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Song With Given ID not found",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: struct {
			SongLikesCount int64 `json:"song_likes_count"`
		}{
			SongLikesCount: s.songRepoService.GetSongLikeCount(uuid.FromStringOrNil(id)),
		},
	}, nil
}

func (s Service) LikeSong(songLike *model2.SongLikes) (*constant.SuccessData, *constant.ErrorData) {
	_, err := s.songRepoService.Song(songLike.SongId)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Song With Given ID not found",
		}
	}
	isUnLiked, err := s.songRepoService.LikeSong(songLike)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error  While Updating Like",
		}
	}
	if isUnLiked {
		return &constant.SuccessData{
			Code: http.StatusNoContent,
			Data: "Album Like Undid",
		}, nil
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Album Liked",
	}, nil
}

func (s Service) CreateSong(song *model2.Song) (*constant.SuccessData, *constant.ErrorData) {
	if song.Title == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Song Title Cannot Be Empty",
		}
	}
	if song.ArtistId.String() == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Artist Id Cannot Be Empty",
		}
	}

	_, err := s.artistRepoService.Artist(song.ArtistId)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Artist With Given ID not found",
		}
	}
	_, err = s.songRepoService.CreateSong(song)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Creating Song",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Song Created",
	}, nil
}

func (s Service) Songs(sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData) {
	isValid := IsValidFilter(sortKey) && (constant.IsValidSort(sort))

	if !isValid {
		sort = ""
		sortKey = ""
	}
	songs, p, err := s.songRepoService.Songs(sort, sortKey, pageNo, maxPerPage)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Getting Songs",
		}
	}
	pageView := view.New(p)

	return &constant.SuccessData{
		Code: 200,
		Data: constant.PaginatedData{
			MetaData: constant.CreateMetaData(
				p.Page(), maxPerPage, p.PageNums(), p.Nums(),
				constant.LinksData{
					Self:     constant.GetFormattedLinkType1(isValid, "v1/songs", pageNo, maxPerPage, sort, sortKey),
					First:    constant.GetFormattedLinkType1(isValid, "v1/songs", 0, maxPerPage, sort, sortKey),
					Previous: constant.GetFormattedLinkType1(isValid, "v1/songs", pageView.Prev(), maxPerPage, sort, sortKey),
					Next:     constant.GetFormattedLinkType1(isValid, "v1/songs", pageView.Next(), maxPerPage, sort, sortKey),
					Last:     constant.GetFormattedLinkType1(isValid, "v1/songs", pageView.Last(), maxPerPage, sort, sortKey),
				}),
			Data: songs,
		},
	}, nil
}

func (s Service) Song(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData) {
	r, err := s.songRepoService.Song(id)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Song With The Given Id not Found",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: r,
	}, nil
}

func (s Service) UpdateSong(song *model2.Song) (*constant.SuccessData, *constant.ErrorData) {
	if song.Title == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Song Title Cannot Be Empty",
		}
	}
	if song.ArtistId.String() == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Artist Id Cannot Be Empty",
		}
	}
	_, err := s.songRepoService.Song(song.ID)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Song With Given ID not found",
		}
	}
	_, err = s.artistRepoService.Artist(song.ArtistId)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Artist With Given ID not found",
		}
	}
	_, err = s.songRepoService.UpdateSong(song)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Updating Song",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Song Updated",
	}, nil
}

func (s Service) DeleteSong(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData) {
	if id.String() == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Album Id Cannot Be Empty",
		}
	}
	_, err := s.songRepoService.Song(id)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Song With The Given ID not Found",
		}
	}
	_, err = s.songRepoService.DeleteSong(id)
	if err != nil {

		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Deleting Song",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Song Deleted",
	}, nil
}

func NewSongService(songRepoService song.SongsService, artistRepoService artist.ArtistsService) UseCase {
	return &Service{
		songRepoService:   songRepoService,
		artistRepoService: artistRepoService,
	}
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
