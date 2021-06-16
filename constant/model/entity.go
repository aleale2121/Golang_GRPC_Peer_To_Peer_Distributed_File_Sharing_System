package model


import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)


type Artist struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"artist_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Image     string    `json:"image"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Artist) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

type ArtistLikes struct {
	ArtistId uuid.UUID `gorm:"type:uuid;" json:"artist_id"`
	UserId   uuid.UUID `gorm:"type:uuid;" json:"user_id"`
}

func (t *ArtistLikes) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

type Favorite struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"favorite_id"`
	Title     string    `json:"title"`
	UserId    uuid.UUID `gorm:"type:uuid;" json:"user_id"`
	SongId    uuid.UUID `gorm:"type:uuid;" json:"song_id"`
	Song      Song      `gorm:"auto_preload" json:"song"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Favorite) BeforeCreate(scope *gorm.Scope) error {
	v4Uuid := uuid.NewV4()

	return scope.SetColumn("ID", v4Uuid)
}
func (t *Favorite) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

type Playlist struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key;" json:"playlist_id"`
	UserId    uuid.UUID       `gorm:"type:uuid;" json:"user_id"`
	Title     string          `json:"title"`
	CreatedBy string          `json:"created_by"`
	Type      string          `json:"type"`
	Songs     []PlaylistSongs `gorm:"ForeignKey:PlaylistId;auto_preload" json:"songs"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (t *Playlist) BeforeCreate(scope *gorm.Scope) error {
	v4Uuid := uuid.NewV4()

	return scope.SetColumn("ID", v4Uuid)
}
func (t *Playlist) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

type PlaylistSongs struct {
	ID         uuid.UUID `json:"id"`
	PlaylistId uuid.UUID `gorm:"type:uuid;" json:"playlist_id"`
	SongId     uuid.UUID `gorm:"type:uuid;" json:"song_id"`
	Song       Song      `gorm:"auto_preload" json:"song"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (t *PlaylistSongs) BeforeCreate(scope *gorm.Scope) error {
	v4Uuid := uuid.NewV4()
	return scope.SetColumn("ID", v4Uuid)
}
func (t *PlaylistSongs) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}


type Song struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;" json:"song_id"`
	ArtistId      uuid.UUID `gorm:"type:uuid;" json:"artist_id"`
	Artist        Artist    `json:"artist"`
	Title         string    `json:"title"`
	CoverImageUrl string    `json:"cover_image_url"`
	SongUrl       string    `json:"song_url"`
	Views         int       `json:"views"`
	Duration      int       `json:"duration"`
	ReleaseAt     time.Time `json:"released_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (s *Song) BeforeCreate(scope *gorm.Scope) error {
	v4Uuid := uuid.NewV4()
	return scope.SetColumn("ID", v4Uuid)
}
func (s *Song) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(s)
}

type SongLikes struct {
	SongId uuid.UUID `gorm:"type:uuid;" json:"song_id"`
	UserId uuid.UUID `gorm:"type:uuid;" json:"user_id"`
}

func (t *SongLikes) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}


type Info struct {
	Id string
	Port int
	Files []string
}
