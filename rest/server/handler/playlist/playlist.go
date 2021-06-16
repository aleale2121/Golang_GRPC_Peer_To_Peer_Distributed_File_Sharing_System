package rest

import (
	"context"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/auth_client_service"
	module "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/module/playlist"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type PlaylistsHandler interface {
	AddPlaylist(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetPlaylists(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetUserPlaylists(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetPlaylistById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeletePlaylist(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdatePlaylist(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	AddSongToPlaylist(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteSongFromPlaylist(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	MiddleWareValidatePlaylist(next httprouter.Handle) httprouter.Handle
	MiddleWareValidatePlaylistSong(next httprouter.Handle) httprouter.Handle
	MiddlewareGetUserId(next httprouter.Handle) httprouter.Handle
}
type playlistsHandler struct {
	playlistsUseCase module.UseCase
}

type keyPlaylist struct{}
type keyPlaylistSong struct{}
type keyUserID struct{}

func NewPlaylistHandler(useCase module.UseCase) PlaylistsHandler {
	return &playlistsHandler{playlistsUseCase: useCase}
}
func (a *playlistsHandler) MiddlewareGetUserId(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		grpcDialStr, err := constant.GetUserManagementGrpcConnectionString()
		if err != nil {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusRequestTimeout,
				Title: "Unable To Create Connection With Remote Server",
			})
			return

		}
		transportOption := grpc.WithInsecure()

		cc, err := grpc.Dial(grpcDialStr, transportOption)

		if err != nil {

			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusRequestTimeout,
				Title: "Unable To Create Connection With Remote Server",
			})
			return
		}
		authClient := *auth_client_service.NewAuthClient(cc)
		token := constant.ExtractToken(r)
		if token == "" {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusUnauthorized,
				Title: "access token is empty",
			})
			return
		}
		userID, err := authClient.GetUserId(token)
		if err != nil || userID == "" {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusUnauthorized,
				Title: err.Error(),
			})
			return
		}

		ctx := context.WithValue(r.Context(), keyUserID{}, userID)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}

func (a playlistsHandler) AddPlaylist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	plList := r.Context().Value(keyPlaylist{}).(model2.Playlist)
	userId := r.Context().Value(keyUserID{}).(string)
	plList.UserId = uuid.FromStringOrNil(userId)
	successData, errData := a.playlistsUseCase.CreatePlaylist(&plList)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a playlistsHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	successData, errData := a.playlistsUseCase.DeletePlaylist(uuid.FromStringOrNil(id))

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a playlistsHandler) MiddleWareValidatePlaylistSong(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		song := model2.PlaylistSongs{}
		err := song.FromJSON(r.Body)
		if err != nil {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusBadRequest,
				Title: "Invalid Playlist Song Data",
			})
			return
		}
		ctx := context.WithValue(r.Context(), keyPlaylistSong{}, song)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}
func (a playlistsHandler) MiddleWareValidatePlaylist(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		playlist := model2.Playlist{}
		err := playlist.FromJSON(r.Body)
		if err != nil {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusBadRequest,
				Title: "Invalid Playlist Data",
			})
			return
		}
		ctx := context.WithValue(r.Context(), keyPlaylist{}, playlist)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}

func (a playlistsHandler) GetPlaylists(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageNo, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		constant.RespondWithError(w, r, constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Unable To get Url parameter for page",
		})
		return
	}
	maxPerPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		constant.RespondWithError(w, r, constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Unable To get Url parameter for per_page",
		})
		return
	}
	sort := r.URL.Query().Get("sort")
	sortKey := r.URL.Query().Get("sort_key")
	successData, errData := a.playlistsUseCase.Playlists(sort, sortKey, pageNo, maxPerPage)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
func (a *playlistsHandler) GetUserPlaylists(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	userId := r.Context().Value(keyUserID{}).(string)
	pageNo, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		constant.RespondWithError(w, r, constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Unable To get Url parameter for page",
		})
		return
	}
	maxPerPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		constant.RespondWithError(w, r, constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Unable To get Url parameter for per_page",
		})
		return
	}
	sort := r.URL.Query().Get("sort")
	sortKey := r.URL.Query().Get("sort_key")
	successData, errData := a.playlistsUseCase.UserPlaylists(userId, sort, sortKey, pageNo, maxPerPage)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
func (a playlistsHandler) GetPlaylistById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	successData, errData := a.playlistsUseCase.Playlist(uuid.FromStringOrNil(id))

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
func (a playlistsHandler) UpdatePlaylist(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	playlist := r.Context().Value(keyPlaylist{}).(model2.Playlist)

	successData, errData := a.playlistsUseCase.UpdatePlaylist(&playlist)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a playlistsHandler) AddSongToPlaylist(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	playlistSong := r.Context().Value(keyPlaylistSong{}).(model2.PlaylistSongs)
	successData, errData := a.playlistsUseCase.AddSongToPlaylist(playlistSong)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a playlistsHandler) DeleteSongFromPlaylist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	successData, errData := a.playlistsUseCase.DeleteSongFromPlaylist(uuid.FromStringOrNil(id))

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
