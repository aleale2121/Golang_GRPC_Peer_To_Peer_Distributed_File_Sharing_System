package playlist

import (
	"fmt"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"time"
)

type playlistsGormRepo struct {
	conn *gorm.DB
}

func (a playlistsGormRepo) GetUserPlaylists(userId string, sort string, sortKey string, pageNo int, maxPerPage int) ([]model2.Playlist, *paginator.Paginator, error) {

	var playlists []model2.Playlist
	var orderBy string
	var q *gorm.DB
	if sort == "" {
		q = a.conn.Preload("Songs").Preload("Songs.Song").Preload("Songs.Song.Artist").Where("user_id=?", userId).Model(model2.Playlist{})
	} else if sort == "ASC" {
		orderBy = fmt.Sprintf("%s ASC", sortKey)
		q = a.conn.Preload("Songs").Preload("Songs.Song").Preload("Songs.Song.Artist").Where("user_id=?", userId).Order(orderBy).Model(model2.Playlist{})
	} else if sort == "DESC" {

		orderBy = fmt.Sprintf("%s DESC", sortKey)
		q = a.conn.Preload("Songs").Preload("Songs.Song").Preload("Songs.Song.Artist").Where("user_id=?", userId).Order(orderBy).Model(model2.Playlist{})
	} else {
		q = a.conn.Preload("Songs").Preload("Songs.Song").Preload("Songs.Song.Artist").Where("user_id=?", userId).Model(model2.Playlist{})

	}

	p := paginator.New(adapter.NewGORMAdapter(q), maxPerPage)
	p.SetPage(pageNo)

	if err := p.Results(&playlists); err != nil {
		return nil, nil, err
	}
	return playlists, &p, nil
}

func (a playlistsGormRepo) IsPlaylistNameExist(title string) bool {
	var count int64
	a.conn.Model(&model2.Playlist{}).Where("title=?", title).Count(&count)
	return count > 0
}

func NewPlaylistsGormRepo(db *gorm.DB) PlaylistsService {
	return &playlistsGormRepo{conn: db}
}

func (a playlistsGormRepo) CreatePlaylist(playlist *model2.Playlist) (*model2.Playlist, error) {
	err := a.conn.Create(playlist).Error
	if err != nil {
		return nil, err
	}
	return playlist, err
}

func (a playlistsGormRepo) Playlists(sort string, sortKey string, pageNo int, maxPerPage int) ([]model2.Playlist, *paginator.Paginator, error) {
	var playlists []model2.Playlist
	var orderBy string
	var q *gorm.DB
	if sort == "" {
		q = a.conn.Preload("Songs").Preload("Songs.Song").Preload("Songs.Song.Artist").Model(model2.Playlist{})
	} else if sort == "ASC" {
		orderBy = fmt.Sprintf("%s ASC", sortKey)
		q = a.conn.Preload("Songs").Preload("Songs.Song").Preload("Songs.Song.Artist").Order(orderBy).Model(model2.Playlist{})
	} else if sort == "DESC" {

		orderBy = fmt.Sprintf("%s DESC", sortKey)
		q = a.conn.Preload("Songs").Preload("Songs.Song").Preload("Songs.Song.Artist").Order(orderBy).Model(model2.Playlist{})
	} else {
		q = a.conn.Preload("Songs").Preload("Songs.Song").Preload("Songs.Song.Artist").Model(model2.Playlist{})

	}

	p := paginator.New(adapter.NewGORMAdapter(q), maxPerPage)
	p.SetPage(pageNo)

	if err := p.Results(&playlists); err != nil {
		return nil, nil, err
	}
	return playlists, &p, nil
}

func (a playlistsGormRepo) Playlist(id uuid.UUID) (*model2.Playlist, error) {
	pl := model2.Playlist{}
	err := a.conn.Preload("Songs").Preload("Songs.Song").Where("id=?", id).First(&pl).Error
	if err != nil {
		return nil, err
	}
	return &pl, err
}

func (a playlistsGormRepo) UpdatePlaylist(playlist *model2.Playlist) (*model2.Playlist, error) {
	updateColumns := a.conn.Preload("Songs").Model(playlist).Where("id = ?", playlist.ID).UpdateColumns(
		map[string]interface{}{
			"title":      playlist.Title,
			"created_by": playlist.CreatedBy,
			"type":       playlist.Type,
			"updated_at": time.Now(),
		},
	)
	if updateColumns.Error != nil {
		return nil, updateColumns.Error
	}
	return playlist, nil
}

func (a playlistsGormRepo) DeletePlaylist(id uuid.UUID) error {
	toBeDeleted, err := a.Playlist(id)
	if err != nil {
		return err
	}
	err = a.conn.Where("id=?", id).Delete(&toBeDeleted).Error
	if err != nil {
		return err
	}
	return nil
}

func (a playlistsGormRepo) AddSongToPlaylist(song *model2.PlaylistSongs) error {
	err := a.conn.Create(song).Error
	if err != nil {
		return err
	}
	return nil
}

func (a playlistsGormRepo) DeleteSongFromPlaylist(id uuid.UUID) error {
	song := model2.PlaylistSongs{}
	err := a.conn.Preload("Songs").Where("id=?", id).First(&song).Error
	if err != nil {
		return err
	}
	err = a.conn.Where("id=?", id).Delete(&song).Error
	if err != nil {
		return err
	}
	return nil
}
