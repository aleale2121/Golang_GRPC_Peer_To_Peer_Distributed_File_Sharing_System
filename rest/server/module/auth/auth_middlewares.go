package auth

import (
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/auth_client_service"
	"google.golang.org/grpc"

	"github.com/julienschmidt/httprouter"
	"net/http"
)

type MiddleWareAuth struct {
	role []string
}

func NewAuthMiddleWare(
	role []string) *MiddleWareAuth {
	return &MiddleWareAuth{
		role: role,
	}

}
func (m *MiddleWareAuth) Authorized(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := constant.ExtractToken(r)
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
		if token == "" {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusUnauthorized,
				Title: "access token is empty",
			})
			return
		}
		isAuthorized, err := authClient.IsAuthorized(token, m.role)
		if err != nil {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusUnauthorized,
				Title: err.Error(),
			})
			return
		}
		if !isAuthorized {
			constant.RespondWithError(w, r, constant.ErrorData{
				Code:  http.StatusUnauthorized,
				Title: "Permission Denied",
			})
			return
		}
		next(w, r, ps)
	}
}
