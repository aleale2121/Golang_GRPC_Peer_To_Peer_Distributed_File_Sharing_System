package favorite

import (
	"fmt"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"time"
)

type favoritesGormRepo struct {
	conn *gorm.DB
}

func (a favoritesGormRepo) GetFavoriteById(favoriteId uuid.UUID) (*model2.Favorite, error) {
	f := model2.Favorite{}
	err := a.conn.Where("id=?", favoriteId).First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, err
}

func NewFavoriteGormRepo(db *gorm.DB) FavoritesService {
	return &favoritesGormRepo{conn: db}
}
func (a favoritesGormRepo) GetUserFavoriteSongs(userId uuid.UUID, sort string, sortKey string, pageNo int, maxPerPage int) ([]model2.Favorite, *paginator.Paginator, error) {

	var favorites []model2.Favorite
	var orderBy string
	var q *gorm.DB
	if sort == "" {
		q = a.conn.Preload("Song").Preload("Song.Artist").Where("user_id=?", userId).Model(model2.Favorite{})

	} else if sort == "ASC" {
		orderBy = fmt.Sprintf("%s ASC", sortKey)
		q = a.conn.Preload("Song").Preload("Song.Artist").Where("user_id=?", userId).Order(orderBy).Model(model2.Favorite{})
	} else if sort == "DESC" {

		orderBy = fmt.Sprintf("%s DESC", sortKey)
		q = a.conn.Preload("Song").Preload("Song.Artist").Where("user_id=?", userId).Order(orderBy).Model(model2.Favorite{})
	} else {
		q = a.conn.Preload("Song").Preload("Song.Artist").Where("user_id=?", userId).Model(model2.Favorite{})
	}
	q = a.conn.Preload("Song").Preload("Song.Artist").Where("user_id=?", userId).Model(model2.Favorite{})
	p := paginator.New(adapter.NewGORMAdapter(q), maxPerPage)
	p.SetPage(pageNo)

	if err := p.Results(&favorites); err != nil {
		return nil, nil, err
	}
	return favorites, &p, nil
}

func (a favoritesGormRepo) GetUserFavoriteSong(userId, songId uuid.UUID) (*model2.Favorite, error) {
	pl := model2.Favorite{}
	err := a.conn.Preload("Song").Preload("Song.Artist").Where("user_id=? AND song_id=?", userId, songId).First(&pl).Error
	if err != nil {
		return nil, err
	}
	return &pl, err
}

func (a favoritesGormRepo) RemoveSongFromUserFavorite(userId, songId uuid.UUID) error {
	toBeDeleted, err := a.GetUserFavoriteSong(userId, songId)
	if err != nil {
		return err
	}
	err = a.conn.Where("user_id=? AND song_id=?", userId, songId).Delete(&toBeDeleted).Error
	if err != nil {
		return err
	}
	return nil
}

func (a favoritesGormRepo) CreateFavorite(favorite *model2.Favorite) (*model2.Favorite, error) {
	err := a.conn.Create(favorite).Error
	if err != nil {
		return nil, err
	}
	return favorite, err
}
func (a favoritesGormRepo) UpdateFavoriteName(favorite model2.Favorite) (*model2.Favorite, error) {
	updateColumns := a.conn.Model(favorite).Where("id = ?", favorite.ID).UpdateColumns(
		map[string]interface{}{
			"title":      favorite.Title,
			"updated_at": time.Now(),
		},
	)
	if updateColumns.Error != nil {
		return nil, updateColumns.Error
	}
	return &favorite, nil
}
