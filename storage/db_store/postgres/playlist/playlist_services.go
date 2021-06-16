package playlist

import (
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator"
)

type PlaylistsService interface {
	CreatePlaylist(song *model2.Playlist) (*model2.Playlist, error)
	Playlists(sort string, sortKey string, pageNo int, maxPerPage int) ([]model2.Playlist, *paginator.Paginator, error)
	GetUserPlaylists(userId string, sort string, sortKey string, pageNo int, maxPerPage int) ([]model2.Playlist, *paginator.Paginator, error)
	Playlist(id uuid.UUID) (*model2.Playlist, error)
	UpdatePlaylist(playlist *model2.Playlist) (*model2.Playlist, error)
	DeletePlaylist(id uuid.UUID) error
	AddSongToPlaylist(song *model2.PlaylistSongs) error
	DeleteSongFromPlaylist(id uuid.UUID) error
	IsPlaylistNameExist(title string) bool
}
