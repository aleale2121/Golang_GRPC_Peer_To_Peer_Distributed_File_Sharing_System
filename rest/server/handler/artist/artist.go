package rest

import (
	"context"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/auth_client_service"
	module "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/module/artist"
	"google.golang.org/grpc"
	"strconv"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"

	"net/http"
)

type ArtistsHandler interface {
	AddArtist(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetArtists(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetArtistById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetAlbumArtistCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteArtist(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateArtist(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	LikeArtist(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	MiddleWareValidateArtist(next httprouter.Handle) httprouter.Handle
	MiddleWareValidateArtistLike(next httprouter.Handle) httprouter.Handle
	MiddlewareGetUserId(next httprouter.Handle) httprouter.Handle
}
type artistsHandler struct {
	artistUseCase module.UseCase
}

func NewArtistHandler(useCase module.UseCase) ArtistsHandler {
	return &artistsHandler{artistUseCase: useCase}
}

type keyArtist struct{}
type keyArtistLike struct{}
type keyUserID struct{}

func (a *artistsHandler) MiddlewareGetUserId(next httprouter.Handle) httprouter.Handle {
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

func (a artistsHandler) MiddleWareValidateArtistLike(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		likes := model2.ArtistLikes{}
		err := likes.FromJSON(r.Body)
		if err != nil {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusBadRequest,
				Title: "Invalid Liking  Data",
			})
			return
		}
		ctx := context.WithValue(r.Context(), keyArtistLike{}, likes)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}
func (a artistsHandler) GetAlbumArtistCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	successData, errData := a.artistUseCase.GetArtistLikeCount(id)
	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
func (a artistsHandler) MiddleWareValidateArtist(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		artistX := model2.Artist{}
		err := artistX.FromJSON(r.Body)
		if err != nil {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusBadRequest,
				Title: "Invalid Artist Data",
			})
			return
		}
		ctx := context.WithValue(r.Context(), keyArtist{}, artistX)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}
func (a artistsHandler) LikeArtist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	like := r.Context().Value(keyArtistLike{}).(model2.ArtistLikes)
	userId := r.Context().Value(keyUserID{}).(string)
	like.UserId = uuid.FromStringOrNil(userId)
	successData, errData := a.artistUseCase.LikeArtist(&like)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
func (a artistsHandler) GetArtists(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	successData, errData := a.artistUseCase.Artists(sort, sortKey, pageNo, maxPerPage)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a artistsHandler) GetArtistById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	successData, errData := a.artistUseCase.Artist(uuid.FromStringOrNil(id))

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a artistsHandler) AddArtist(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ar := r.Context().Value(keyArtist{}).(model2.Artist)
	successData, errData := a.artistUseCase.CreateArtist(&ar)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a artistsHandler) DeleteArtist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	successData, errData := a.artistUseCase.DeleteArtist(uuid.FromStringOrNil(id))

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a artistsHandler) UpdateArtist(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ar := r.Context().Value(keyArtist{}).(model2.Artist)

	successData, errData := a.artistUseCase.UpdateArtist(&ar)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
