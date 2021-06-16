package song

import (
	"github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator"
)

type SongsService interface {
	CreateSong(artist *model.Song) (*model.Song, error)
	Songs(sort string, sortKey string, pageNo int, maxPerPage int) ([]model.Song, *paginator.Paginator, error)
	SongByArtistID(artistId string, sort string, sortKey string, pageNo int, maxPerPage int) ([]model.Song, *paginator.Paginator, error)
	Song(id uuid.UUID) (*model.Song, error)
	UpdateSong(artist *model.Song) (*model.Song, error)
	DeleteSong(id uuid.UUID) (*model.Song, error)
	LikeSong(songLike *model.SongLikes) (bool, error)
	GetSongLikeCount(id uuid.UUID) int64
	IncreaseSongViews(id string) error
	GetSongViewsCount(id uuid.UUID) (int64, error)
}
