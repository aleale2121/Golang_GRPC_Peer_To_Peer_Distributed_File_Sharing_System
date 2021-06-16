package rest

import (
	"context"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	model2 "github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/auth_client_service"
	module "github.com/aleale2121/DSP_LAB/Music_Service/rest/server/module/favorite"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type FavoriteHandler interface {
	AddSongToFavorite(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetUserFavoriteSongs(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateFavorite(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	RemoveSongFromUserFavorites(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	MiddleWareValidateFavorite(next httprouter.Handle) httprouter.Handle
	MiddlewareGetUserId(next httprouter.Handle) httprouter.Handle
}
type favoriteHandler struct {
	favoriteUseCase module.UseCase
}

func NewFavoriteHandler(favUseCase module.UseCase) FavoriteHandler {
	return &favoriteHandler{favoriteUseCase: favUseCase}
}

type keyFavorite struct{}
type keyUserID struct{}

func (a *favoriteHandler) MiddlewareGetUserId(next httprouter.Handle) httprouter.Handle {
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

func (a favoriteHandler) MiddleWareValidateFavorite(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fav := model2.Favorite{}
		err := fav.FromJSON(r.Body)
		if err != nil {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusBadRequest,
				Title: "Invalid Favorite Data",
			})
			return
		}
		ctx := context.WithValue(r.Context(), keyFavorite{}, fav)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}

func (a favoriteHandler) AddSongToFavorite(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fav := r.Context().Value(keyFavorite{}).(model2.Favorite)
	userId := r.Context().Value(keyUserID{}).(string)
	fav.UserId = uuid.FromStringOrNil(userId)
	successData, errData := a.favoriteUseCase.CreateFavorite(&fav)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a favoriteHandler) GetUserFavoriteSongs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	successData, errData := a.favoriteUseCase.GetUserFavoriteSongs(uuid.FromStringOrNil(userId), sort, sortKey, pageNo, maxPerPage)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}

func (a favoriteHandler) UpdateFavorite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	favId := ps.ByName("favId")
	title := r.URL.Query().Get("title")

	successData, errData := a.favoriteUseCase.UpdateFavorite(title, favId)

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return

}
func (a favoriteHandler) RemoveSongFromUserFavorites(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := r.Context().Value(keyUserID{}).(string)
	songId := ps.ByName("songId")
	successData, errData := a.favoriteUseCase.RemoveSongFromUserFavorite(uuid.FromStringOrNil(userId), uuid.FromStringOrNil(songId))

	if errData != nil {
		constant.RespondWithError(w, r, *errData)
		return
	}
	constant.RespondWithSuccess(w, *successData)
	return
}
