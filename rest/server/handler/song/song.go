package rest

import (
	"context"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/auth_client_service"
	module "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/module/song"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type SongsHandler interface {
	AddSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetSongs(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetArtistSongs(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetSongById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetSongLikeCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	LikeSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	IncreaseSongViews(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetSongViewsCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	MiddleWareValidateSong(next httprouter.Handle) httprouter.Handle
	MiddleWareValidateSongLike(next httprouter.Handle) httprouter.Handle
	MiddlewareGetUserId(next httprouter.Handle) httprouter.Handle
}
type songsHandler struct {
	songUseCase module.UseCase
}
type keyUserID struct{}

func (a *songsHandler) MiddlewareGetUserId(next httprouter.Handle) httprouter.Handle {
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
func (a songsHandler) IncreaseSongViews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	successData, errData := a.songUseCase.IncreaseSongViews(id)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a songsHandler) GetSongViewsCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	successData, errData := a.songUseCase.GetSongViewsCount(id)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func NewSongHandler(useCase module.UseCase) SongsHandler {
	return &songsHandler{songUseCase: useCase}
}

type keySong struct{}
type keySongLike struct{}

func (a songsHandler) MiddleWareValidateSongLike(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		likes := model2.SongLikes{}
		err := likes.FromJSON(r.Body)
		if err != nil {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusBadRequest,
				Title: "Invalid Liking  Data",
			})
			return
		}
		ctx := context.WithValue(r.Context(), keySongLike{}, likes)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}
func (a songsHandler) LikeSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	like := r.Context().Value(keySongLike{}).(model2.SongLikes)
	userId := r.Context().Value(keyUserID{}).(string)
	like.UserId = uuid.FromStringOrNil(userId)
	successData, errData := a.songUseCase.LikeSong(&like)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
func (a songsHandler) GetSongLikeCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	successData, errData := a.songUseCase.GetSongLikeCount(id)
	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a songsHandler) MiddleWareValidateSong(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		songX := model2.Song{}
		err := songX.FromJSON(r.Body)
		if err != nil {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusBadRequest,
				Title: "Invalid Song Data",
			})
			return
		}
		ctx := context.WithValue(r.Context(), keySong{}, songX)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}

func (a songsHandler) GetArtistSongs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	artistID := r.URL.Query().Get("artist_id")
	if artistID == "" {
		constant.RespondWithError(w, r, constant.ErrorData{
			Code:  http.StatusBadRequest,
			Title: "Unable To get Url parameter for artist ID",
		})
		return
	}
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
	successData, errData := a.songUseCase.ArtistSongs(artistID, sort, sortKey, pageNo, maxPerPage)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
func (a songsHandler) GetSongs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	successData, errData := a.songUseCase.Songs(sort, sortKey, pageNo, maxPerPage)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a songsHandler) GetSongById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	successData, errData := a.songUseCase.Song(uuid.FromStringOrNil(id))

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a songsHandler) AddSong(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sg := r.Context().Value(keySong{}).(model2.Song)
	successData, errData := a.songUseCase.CreateSong(&sg)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a songsHandler) DeleteSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	successData, errData := a.songUseCase.DeleteSong(uuid.FromStringOrNil(id))

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a songsHandler) UpdateSong(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sg := r.Context().Value(keySong{}).(model2.Song)

	successData, errData := a.songUseCase.UpdateSong(&sg)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
