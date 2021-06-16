package module

import (
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/playlist"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/song"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator/view"
	"net/http"
)

type UseCase interface {
	CreatePlaylist(playlist *model2.Playlist) (*constant.SuccessData, *constant.ErrorData)
	Playlists(sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData)
	UserPlaylists(userID string, sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData)
	Playlist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData)
	UpdatePlaylist(playlist *model2.Playlist) (*constant.SuccessData, *constant.ErrorData)
	DeletePlaylist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData)
	AddSongToPlaylist(song model2.PlaylistSongs) (*constant.SuccessData, *constant.ErrorData)
	DeleteSongFromPlaylist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData)
}
type Service struct {
	songRepoService     song.SongsService
	playlistRepoService playlist.PlaylistsService
}

func (s Service) UserPlaylists(userID string, sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData) {

	isValid := IsValidFilter(sortKey) && (constant.IsValidSort(sort))
	if !isValid {
		sort = ""
		sortKey = ""
	}
	playlists, p, err := s.playlistRepoService.GetUserPlaylists(userID, sort, sortKey, pageNo, maxPerPage)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Getting Playlists",
		}
	}
	pageView := view.New(p)

	return &constant.SuccessData{
		Code: 200,
		Data: constant.PaginatedData{
			MetaData: constant.CreateMetaData(
				p.Page(), maxPerPage, p.PageNums(), p.Nums(),
				constant.LinksData{
					Self:     constant.GetFormattedLinkType1(isValid, "v1/user/playlists", pageNo, maxPerPage, sort, sortKey),
					First:    constant.GetFormattedLinkType1(isValid, "v1/user/playlists", 0, maxPerPage, sort, sortKey),
					Previous: constant.GetFormattedLinkType1(isValid, "v1/user/playlists", pageView.Prev(), maxPerPage, sort, sortKey),
					Next:     constant.GetFormattedLinkType1(isValid, "v1/user/playlists", pageView.Next(), maxPerPage, sort, sortKey),
					Last:     constant.GetFormattedLinkType1(isValid, "v1/user/playlists", pageView.Last(), maxPerPage, sort, sortKey),
				}),
			Data: playlists,
		},
	}, nil
}

func (s Service) CreatePlaylist(playlist *model2.Playlist) (*constant.SuccessData, *constant.ErrorData) {
	if playlist.Title == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Playlist Title Cannot Be Empty",
		}
	}
	if playlist.Type == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Playlist Type Cannot Be Empty",
		}
	}
	if s.playlistRepoService.IsPlaylistNameExist(playlist.Title) {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Playlist Name Already Exist",
		}
	}

	_, err := s.playlistRepoService.Playlist(playlist.ID)
	if err == nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Playlist Already Created",
		}
	}
	_, err = s.playlistRepoService.CreatePlaylist(playlist)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Creating Playlist",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Playlist Created",
	}, nil
}

func (s Service) Playlists(sort string, sortKey string, pageNo int, maxPerPage int) (*constant.SuccessData, *constant.ErrorData) {
	isValid := IsValidFilter(sortKey) && (constant.IsValidSort(sort))
	if !isValid {
		sort = ""
		sortKey = ""
	}
	playlists, p, err := s.playlistRepoService.Playlists(sort, sortKey, pageNo, maxPerPage)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Getting Playlists",
		}
	}
	pageView := view.New(p)

	return &constant.SuccessData{
		Code: 200,
		Data: constant.PaginatedData{
			MetaData: constant.CreateMetaData(
				p.Page(), maxPerPage, p.PageNums(), p.Nums(),
				constant.LinksData{
					Self:     constant.GetFormattedLinkType1(isValid, "v1/playlists", pageNo, maxPerPage, sort, sortKey),
					First:    constant.GetFormattedLinkType1(isValid, "v1/playlists", 0, maxPerPage, sort, sortKey),
					Previous: constant.GetFormattedLinkType1(isValid, "v1/playlists", pageView.Prev(), maxPerPage, sort, sortKey),
					Next:     constant.GetFormattedLinkType1(isValid, "v1/playlists", pageView.Next(), maxPerPage, sort, sortKey),
					Last:     constant.GetFormattedLinkType1(isValid, "v1/playlists", pageView.Last(), maxPerPage, sort, sortKey),
				}),
			Data: playlists,
		},
	}, nil
}

func (s Service) Playlist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData) {
	r, err := s.playlistRepoService.Playlist(id)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Playlist With The Given Id not Found",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: r,
	}, nil
}

func (s Service) UpdatePlaylist(playlist *model2.Playlist) (*constant.SuccessData, *constant.ErrorData) {
	if playlist.Title == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Playlist Name Cannot Be Empty",
		}
	}
	if playlist.Type == "" {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Playlist Type Cannot Be Empty",
		}
	}
	if s.playlistRepoService.IsPlaylistNameExist(playlist.Title) {
		return nil, &constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Playlist Name Already Exist",
		}
	}
	_, err := s.playlistRepoService.Playlist(playlist.ID)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Playlist With Given  ID not found",
		}
	}

	_, err = s.playlistRepoService.UpdatePlaylist(playlist)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Updating Playlist",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Playlist Updated",
	}, nil
}

func (s Service) DeletePlaylist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData) {

	_, err := s.playlistRepoService.Playlist(id)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Playlist With The Given ID not Found",
		}
	}
	err = s.playlistRepoService.DeletePlaylist(id)
	if err != nil {

		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Deleting Album",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Playlist Deleted",
	}, nil
}

func (s Service) AddSongToPlaylist(song model2.PlaylistSongs) (*constant.SuccessData, *constant.ErrorData) {
	_, err := s.playlistRepoService.Playlist(song.PlaylistId)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "No Playlist Found With The Given ID",
		}
	}
	_, err = s.songRepoService.Song(song.SongId)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "No Song Found With The Given ID",
		}
	}
	err = s.playlistRepoService.AddSongToPlaylist(&song)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Adding Song To Playlist",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Song Added To Playlist",
	}, nil
}

func (s Service) DeleteSongFromPlaylist(id uuid.UUID) (*constant.SuccessData, *constant.ErrorData) {
	err := s.playlistRepoService.DeleteSongFromPlaylist(id)
	if err != nil {
		return nil, &constant.ErrorData{
			Code:  http.StatusNotFound,
			Title: "Song With The Given ID not Found In The Playlist",
		}
	}
	err = s.playlistRepoService.DeletePlaylist(id)
	if err != nil {

		return nil, &constant.ErrorData{
			Code:  http.StatusInternalServerError,
			Title: "Error Occurred While Deleting Album",
		}
	}
	return &constant.SuccessData{
		Code: 200,
		Data: "Playlist Deleted",
	}, nil
}

func NewPlaylistService(songRepoService song.SongsService, playlistRepoService playlist.PlaylistsService) UseCase {
	return &Service{
		playlistRepoService: playlistRepoService,
		songRepoService:     songRepoService,
	}
}
func IsValidFilter(key string) bool {
	keys := [3]string{"title", "created_by", "type"}
	for i := 0; i < len(keys); i++ {
		if keys[i] == key {
			return true
		}
	}
	return false
}
