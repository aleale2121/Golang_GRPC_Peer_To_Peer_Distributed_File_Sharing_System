package main

import (
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	ah "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/handler/artist"
	fh "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/handler/favorite"
	ph "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/handler/playlist"
	sh "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/handler/song"
	arm "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/module/artist"
	"github.com/aleale2121/DSP_LAB/Music_Service/rest/server/module/auth"
	fm "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/module/favorite"
	pm "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/module/playlist"
	sm "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/module/song"
	ar "github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/artist"
	fr "github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/favorite"
	pr "github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/playlist"
	sr "github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/song"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

func main() {

	connStr, dialect, err := constant.GetGormDatabaseConnectionString()
	if err != nil {
		panic(err)
	}
	dbConn, err := gorm.Open(dialect,
		connStr)
	if dbConn != nil {
		defer dbConn.Close()
	}
	if err != nil {
		panic(err)
	}
	httpRouter := httprouter.New()

	artistRepoService := ar.NewArtistGormRepo(dbConn)
	artistUseCase := arm.NewArtistService(artistRepoService)
	aH := ah.NewArtistHandler(artistUseCase)

	adminOnly := auth.NewAuthMiddleWare([]string{"ADMIN"})
	adminAndArtist := auth.NewAuthMiddleWare([]string{"ADMIN", "ARTIST"})
	all := auth.NewAuthMiddleWare([]string{"ADMIN", "USER", "ARTIST"})

	httpRouter.Handle(http.MethodGet,"/v1/artists", all.Authorized(aH.GetArtists))
	httpRouter.Handle(http.MethodGet,"/v1/artists/:id", all.Authorized(aH.GetArtistById))
	httpRouter.Handle(http.MethodPost,"/v1/artists",adminOnly.Authorized(aH.MiddleWareValidateArtist(aH.AddArtist)))
	httpRouter.Handle(http.MethodPut,"/v1/artists",adminAndArtist.Authorized(aH.MiddleWareValidateArtist(aH.UpdateArtist)))
	httpRouter.Handle(http.MethodDelete,"/v1/artists/:id",adminAndArtist.Authorized(aH.DeleteArtist))
	httpRouter.Handle(http.MethodPost,"/v1/like/artist",all.Authorized(
		aH.MiddlewareGetUserId(aH.MiddleWareValidateArtistLike(aH.LikeArtist))))
	httpRouter.Handle(http.MethodGet,"/v1/like/artist/:id",all.Authorized(
		aH.MiddlewareGetUserId(aH.MiddleWareValidateArtistLike(aH.GetAlbumArtistCount))))

	songRepoService := sr.NewSongGormRepo(dbConn)
	songUseCase := sm.NewSongService(songRepoService, artistRepoService)
	sgH := sh.NewSongHandler(songUseCase)
	httpRouter.Handle(http.MethodGet,"/v1/songs", all.Authorized(sgH.GetSongs))
	httpRouter.Handle(http.MethodGet,"/v1/songs/:id", all.Authorized(sgH.GetSongById))
	httpRouter.Handle(http.MethodPost,"/v1/songs",adminAndArtist.Authorized(sgH.MiddleWareValidateSong(sgH.AddSong)))
	httpRouter.Handle(http.MethodPut,"/v1/songs",adminAndArtist.Authorized(sgH.MiddleWareValidateSong(sgH.UpdateSong)))
	httpRouter.Handle(http.MethodDelete,"/v1/songs/:id",adminAndArtist.Authorized(sgH.DeleteSong))
	httpRouter.Handle(http.MethodPost,"/v1/like/song",all.Authorized(
		sgH.MiddlewareGetUserId(sgH.MiddleWareValidateSongLike(sgH.LikeSong))))
	httpRouter.Handle(http.MethodGet,"/v1/like/song/:id",all.Authorized(
		sgH.GetSongLikeCount))
	httpRouter.Handle(http.MethodPost,"/v1/view/song/:id",all.Authorized(
		sgH.IncreaseSongViews))
	httpRouter.Handle(http.MethodGet,"/v1/view/song/:id",all.Authorized(
		sgH.GetSongViewsCount))
	httpRouter.Handle(http.MethodGet,"/v1/artist/songs", all.Authorized(sgH.GetArtistSongs))




	playlistRepoService := pr.NewPlaylistsGormRepo(dbConn)
	playlistUseCase := pm.NewPlaylistService(songRepoService, playlistRepoService)
	plH := ph.NewPlaylistHandler(playlistUseCase)

	httpRouter.Handle(http.MethodGet,"/v1/playlists", all.Authorized(plH.GetPlaylists))
	httpRouter.Handle(http.MethodGet,"/v1/playlists/:id", all.Authorized(plH.GetPlaylistById))
	httpRouter.Handle(http.MethodPost,"/v1/playlists",adminAndArtist.Authorized(plH.MiddlewareGetUserId(plH.MiddleWareValidatePlaylist(plH.AddPlaylist))))
	httpRouter.Handle(http.MethodPut,"/v1/playlists",adminAndArtist.Authorized(plH.MiddleWareValidatePlaylist(plH.UpdatePlaylist)))
	httpRouter.Handle(http.MethodDelete,"/v1/playlists/:id",adminAndArtist.Authorized(plH.DeletePlaylist))
	httpRouter.Handle(http.MethodPost,"/v1/playlist/song",adminAndArtist.Authorized(plH.MiddleWareValidatePlaylistSong(plH.AddSongToPlaylist)))
	httpRouter.Handle(http.MethodGet,"/v1/user/playlists",all.Authorized(plH.MiddlewareGetUserId(plH.GetUserPlaylists)))
	httpRouter.Handle(http.MethodDelete,"/v1/playlist/song/:id",adminAndArtist.Authorized(plH.DeleteSongFromPlaylist))


	favoriteRepoService := fr.NewFavoriteGormRepo(dbConn)
	favoriteUseCase := fm.NewFavoriteService(favoriteRepoService)
	fvH:= fh.NewFavoriteHandler(favoriteUseCase)
	httpRouter.Handle(http.MethodGet,"/v1/favorites/:userId", all.Authorized(fvH.MiddlewareGetUserId(fvH.GetUserFavoriteSongs)))
	httpRouter.Handle(http.MethodPost,"/v1/favorites",all.Authorized(fvH.MiddlewareGetUserId(fvH.AddSongToFavorite)))
	httpRouter.Handle(http.MethodPut,"/v1/favorites/:favId",all.Authorized(fvH.UpdateFavorite))
	httpRouter.Handle(http.MethodDelete,"/v1/favorites/:id",all.Authorized(fvH.MiddlewareGetUserId(fvH.RemoveSongFromUserFavorites)))

	curDir, err:= os.Getwd()
	if err != nil {
		panic(err)
	}
	httpRouter.ServeFiles("/assets/images/*filepath", http.Dir(curDir+"/assets/images/"))
	err = http.ListenAndServe("localhost:8080", &Server{httpRouter})
	if err != nil {
		panic(err)
	}
}

type Server struct {
	r *httprouter.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	s.r.ServeHTTP(w, r)
}
