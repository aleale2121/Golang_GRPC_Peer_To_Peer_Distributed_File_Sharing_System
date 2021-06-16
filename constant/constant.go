package constant

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"strings"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
type ErrorResponse struct {
	Success bool         `json:"success"`
	Errors  ErrorMessage `json:"errors"`
}
type ErrorMessage struct {
	Code   int    `json:"code"`
	Source string `json:"source"`
	Title  string `json:"title"`
}
type ErrorData struct {
	Code  int
	Title string
}
type SuccessData struct {
	Code int
	Data interface{}
}

type PaginatedData struct {
	MetaData MetaData    `json:"meta_data"`
	Data     interface{} `json:"data"`
}
type MetaData struct {
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	PageCount  int                 `json:"page_count"`
	TotalCount int                 `json:"total_count"`
	Links      []map[string]string `json:"links"`
}

func CreateMetaData(page int, perPage int, pageCount int, totalCount int, data LinksData) MetaData {
	return MetaData{
		Page:       page,
		PerPage:    perPage,
		PageCount:  pageCount,
		TotalCount: totalCount,
		Links: []map[string]string{
			{"self": data.Self},
			{"first": data.First},
			{"previous": data.Previous},
			{"next": data.Next},
			{"last": data.Last},
		},
	}
}

func GetFormattedLinkType1(isValid bool, name string, page int, maxPerPage int, sort string, sortKey string) string {
	if isValid {

		return fmt.Sprintf("/%s/?page=%d&per_page=%d&sort=%s&sort_key=%s", name, page, maxPerPage, sort, sortKey)
	}
	return fmt.Sprintf("/%s/?page=%d&per_page=%d", name, page, maxPerPage)

}
func GetFormattedLinkType2(isValid bool, name string, page int, maxPerPage int, sort string, sortKey string) string {
	if isValid {
		return fmt.Sprintf("/%s&?page=%d&per_page=%d&sort=%s&sort_key=%s", name, page, maxPerPage, sort, sortKey)
	}
	return fmt.Sprintf("/%s&?page=%d&per_page=%d", name, page, maxPerPage)
}

type LinksData struct {
	Self     string
	First    string
	Previous string
	Next     string
	Last     string
}

func RespondWithError(w http.ResponseWriter, r *http.Request, error ErrorData) {
	w.Header().Set("Content-type", "application/json; charset-UTF8")
	//statusCode, message := GetStatusCode(error)
	w.WriteHeader(error.Code)
	errResp := ErrorResponse{
		Success: false,
		Errors: ErrorMessage{
			Code:   error.Code,
			Source: r.RequestURI,
			Title:  error.Title,
		},
	}
	output, _ := json.MarshalIndent(errResp, "", "\t\t")
	_, _ = w.Write(output)
}

func RespondWithSuccess(w http.ResponseWriter, successData SuccessData) {
	w.Header().Set("Content-type", "application/json; charset-UTF8")
	w.WriteHeader(successData.Code)
	resp := SuccessResponse{
		Success: true,
		Data:    successData.Data,
	}
	output, _ := json.MarshalIndent(resp, "", "\t\t")
	_, _ = w.Write(output)
}
func IsValidSort(key string) bool {
	keys := [4]string{"ASC", "DESC", "asc", "desc"}
	for i := 0; i < len(keys); i++ {
		if keys[i] == key {
			return true
		}
	}
	return false
}

func GetStatusCode(err error) (int, string) {
	st, ok := status.FromError(err)
	if ok {
		if st.Code() == codes.Canceled {
			return http.StatusNotImplemented, st.Message()

		} else if st.Code() == codes.InvalidArgument {
			return http.StatusBadRequest, "Invalid Inputs "
		} else if st.Code() == codes.DeadlineExceeded {
			return http.StatusRequestTimeout, "Request Timeout"
		} else if st.Code() == codes.NotFound {
			return http.StatusNotFound, st.Message()
		} else if st.Code() == codes.AlreadyExists {
			return http.StatusBadRequest, st.Message()
		} else if st.Code() == codes.PermissionDenied {
			return http.StatusUnauthorized, st.Message()
		} else if st.Code() == codes.Unauthenticated {
			return http.StatusUnauthorized, st.Message()
		} else {
			return http.StatusInternalServerError, st.Message()
		}
	}
	return http.StatusInternalServerError, "Unknown Error Occurred"
}
func GetGormDatabaseConnectionString() (string, string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME")),
		os.Getenv("DB_USER"), nil

}
func GetGrpcConnectionString() (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", err
	}
	return os.Getenv("GPRC_CONN_STRING"), nil
}
func GetUserManagementGrpcConnectionString() (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", err
	}
	return os.Getenv("USER_MANAGEMENT_GPRC_DIAL_STRING"), nil
}
func GetHostNameString() (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", err
	}
	return os.Getenv("HOST_NAME"), nil
}
func GetHost() (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", err
	}
	return os.Getenv("HOST"), nil
}
func GetHostPort() (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", err
	}
	return os.Getenv("HOST_PORT"), nil
}
func GetGrpcDialString() (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", err
	}
	return os.Getenv("GPRC_DIAL_STRING"), nil
}
func GetTestUser() (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", err
	}
	return os.Getenv("TEST_USER"), nil
}
func GetTestPassword() (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", err
	}
	return os.Getenv("TEST_PASSWORD"), nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("auth")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
