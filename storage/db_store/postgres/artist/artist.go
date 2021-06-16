package artist

import (
	"fmt"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"time"
)

type ArtistsGormRepo struct {
	conn *gorm.DB
}

func NewArtistGormRepo(db *gorm.DB) ArtistsService {
	return &ArtistsGormRepo{conn: db}
}
func (a ArtistsGormRepo) GetArtistLikeCount(id uuid.UUID) int64 {
	var count int64
	a.conn.Model(&model2.ArtistLikes{}).Where("artist_id=?", id).Count(&count)
	return count

}

func (a ArtistsGormRepo) LikeArtist(artistLike *model2.ArtistLikes) (bool, error) {
	var count int
	a.conn.Model(&model2.ArtistLikes{}).Where("user_id=? AND artist_id=?", artistLike.UserId, artistLike.ArtistId).Count(&count)
	if count > 0 {
		err := a.conn.Where("user_id=? AND artist_id=?", artistLike.UserId, artistLike.ArtistId).Delete(&artistLike).Error
		if err != nil {
			return false, err
		}
		return true, nil
	}
	err := a.conn.Create(artistLike).Error
	if err != nil {
		return false, err
	}
	return false, nil
}

func (a ArtistsGormRepo) CreateArtist(artist *model2.Artist) (*model2.Artist, error) {
	err := a.conn.Create(artist).Error
	if err != nil {
		return nil, err
	}
	return artist, err
}
func (a ArtistsGormRepo) Artists(sort string, sortKey string, pageNo int, maxPerPage int) ([]model2.Artist, *paginator.Paginator, error) {
	var artists []model2.Artist
	var orderBy string
	var q *gorm.DB
	if sort == "" {
		q = a.conn.Model(model2.Artist{})
	} else if sort == "ASC" {
		orderBy = fmt.Sprintf("%s ASC", sortKey)
		q = a.conn.Order(orderBy).Model(model2.Artist{})
	} else if sort == "DESC" {

		orderBy = fmt.Sprintf("%s DESC", sortKey)
		q = a.conn.Order(orderBy).Model(model2.Artist{})
	} else {
		q = a.conn.Model(model2.Artist{})

	}

	p := paginator.New(adapter.NewGORMAdapter(q), maxPerPage)
	p.SetPage(pageNo)

	if err := p.Results(&artists); err != nil {
		return nil, nil, err
	}
	return artists, &p, nil
}

func (a ArtistsGormRepo) Artist(id uuid.UUID) (*model2.Artist, error) {
	post := model2.Artist{}
	err := a.conn.Where("id=?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, err
}

func (a ArtistsGormRepo) UpdateArtist(artist *model2.Artist) (*model2.Artist, error) {
	fmt.Println("--------updating----")
	fmt.Println(artist)
	updateColumns := a.conn.Model(artist).Where("id = ?", artist.ID).UpdateColumns(
		map[string]interface{}{
			"first_name": artist.FirstName,
			"last_name":  artist.LastName,
			"image":      artist.Image,
			"email":      artist.Email,
			"updated_at": time.Now(),
		},
	)
	if updateColumns.Error != nil {
		return nil, updateColumns.Error
	}
	return artist, nil
}

func (a ArtistsGormRepo) DeleteArtist(id uuid.UUID) (*model2.Artist, error) {
	albumToBeDeleted, err := a.Artist(id)
	if err != nil {
		return nil, err
	}
	err = a.conn.Where("id=?", id).Delete(&albumToBeDeleted).Error
	if err != nil {
		return nil, err
	}
	return albumToBeDeleted, err
}
