package song

import (
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"time"
)

type SongsGormRepo struct {
	conn *gorm.DB
}

func (a SongsGormRepo) IncreaseSongViews(id string) error {
	sg := model.Song{}
	err := a.conn.Where("id=?", uuid.FromStringOrNil(id)).First(&sg).Error
	if err != nil {
		return err
	}
	updateColumns := a.conn.Model(&model.Song{}).Where("id = ?", uuid.FromStringOrNil(id)).Update("views", sg.Views+1)
	if updateColumns.Error != nil {
		return updateColumns.Error
	}
	return nil
}

func (a SongsGormRepo) GetSongViewsCount(id uuid.UUID) (int64, error) {
	song := model.Song{}
	err := a.conn.Where("id=?", id).First(&song).Error
	if err != nil {
		return 0, err
	}
	return int64(song.Views), nil
}

func (a SongsGormRepo) SongByArtistID(artistId string, sort string, sortKey string, pageNo int, maxPerPage int) ([]model.Song, *paginator.Paginator, error) {

	var songs []model.Song
	var orderBy string
	var q *gorm.DB
	if sort == "" {
		q = a.conn.Preload("Artist").Model(model.Song{}).Where("artist_id=?", uuid.FromStringOrNil(artistId))

	} else if sort == "ASC" {
		orderBy = fmt.Sprintf("%s ASC", sortKey)
		q = a.conn.Preload("Artist").Order(orderBy).Model(model.Song{}).Where("artist_id=?", uuid.FromStringOrNil(artistId))
	} else if sort == "DESC" {

		orderBy = fmt.Sprintf("%s DESC", sortKey)
		q = a.conn.Preload("Artist").Order(orderBy).Model(model.Song{}).Where("artist_id=?", uuid.FromStringOrNil(artistId))
	} else {
		q = a.conn.Preload("Artist").Model(model.Song{}).Where("artist_id=?", uuid.FromStringOrNil(artistId))

	}
	p := paginator.New(adapter.NewGORMAdapter(q), maxPerPage)
	p.SetPage(pageNo)

	if err := p.Results(&songs); err != nil {
		return nil, nil, err
	}
	return songs, &p, nil
}

func (a SongsGormRepo) LikeSong(songLike *model.SongLikes) (bool, error) {
	var count int
	a.conn.Model(&model.SongLikes{}).Where("user_id=? AND song_id=?", songLike.UserId, songLike.SongId).Count(&count)
	if count > 0 {
		err := a.conn.Where("user_id=? AND song_id=?", songLike.UserId, songLike.SongId).Delete(&songLike).Error
		if err != nil {
			return false, err
		}
		return true, nil
	}
	err := a.conn.Create(songLike).Error
	if err != nil {
		return false, err
	}
	return false, nil
}

func (a SongsGormRepo) GetSongLikeCount(id uuid.UUID) int64 {
	var count int64
	a.conn.Model(&model.SongLikes{}).Where("song_id=?", id).Count(&count)
	return count
}

func (a SongsGormRepo) CreateSong(song *model.Song) (*model.Song, error) {
	err := a.conn.Create(song).Error
	if err != nil {
		return nil, err
	}
	return song, err
}

func (a SongsGormRepo) Songs(sort string, sortKey string, pageNo int, maxPerPage int) ([]model.Song, *paginator.Paginator, error) {
	var songs []model.Song
	var orderBy string
	var q *gorm.DB
	if sort == "" {
		q = a.conn.Preload("Artist").Model(model.Song{})
	} else if sort == "ASC" {
		orderBy = fmt.Sprintf("%s ASC", sortKey)
		q = a.conn.Preload("Artist").Order(orderBy).Model(model.Song{})
	} else if sort == "DESC" {

		orderBy = fmt.Sprintf("%s DESC", sortKey)
		q = a.conn.Preload("Artist").Order(orderBy).Model(model.Song{})
	} else {
		q = a.conn.Preload("Artist").Model(model.Song{})
	}
	p := paginator.New(adapter.NewGORMAdapter(q), maxPerPage)
	p.SetPage(pageNo)

	if err := p.Results(&songs); err != nil {
		return nil, nil, err
	}
	return songs, &p, nil
}

func (a SongsGormRepo) Song(id uuid.UUID) (*model.Song, error) {
	sg := model.Song{}
	err := a.conn.Preload("Artist").Where("id=?", id).First(&sg).Error
	if err != nil {
		return nil, err
	}
	return &sg, err
}

func (a SongsGormRepo) UpdateSong(song *model.Song) (*model.Song, error) {

	updateColumns := a.conn.Preload("Artist").Model(song).Where("id = ?", song.ID).UpdateColumns(
		map[string]interface{}{
			"artist_id":       song.ArtistId,
			"title":           song.Title,
			"cover_image_url": song.CoverImageUrl,
			"song_url":        song.SongUrl,
			"views":           song.Views,
			"duration":        song.Duration,
			"release_at":      time.Now(),
			"updated_at":      time.Now(),
		},
	)
	if updateColumns.Error != nil {
		return nil, updateColumns.Error
	}
	return song, nil
}

func (a SongsGormRepo) DeleteSong(id uuid.UUID) (*model.Song, error) {
	toBeDeleted, err := a.Song(id)
	if err != nil {
		return nil, err
	}
	err = a.conn.Preload("Artist").Where("id=?", id).Delete(&toBeDeleted).Error
	if err != nil {
		return nil, err
	}
	return toBeDeleted, err
}

func NewSongGormRepo(db *gorm.DB) SongsService {
	return &SongsGormRepo{conn: db}
}
