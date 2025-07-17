package http

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	userServiceGrpc "github.com/gorkagg10/lovify/lovify-user-service/grpc/user-service"
)

type CreateUserResponse struct {
	UserID string `json:"id"`
}

func (h *Handler) LoginSpotify(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	response, err := h.UserServiceClient.MusicProviderLogin(
		r.Context(),
		&userServiceGrpc.MusicProviderLoginRequest{
			UserID: &userID,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, response.GetUrl(), http.StatusTemporaryRedirect)
}

func (h *Handler) RegisterSpotify(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	_, err := h.UserServiceClient.MusicProviderOAuthCallback(
		r.Context(),
		&userServiceGrpc.MusicProviderOAuthCallbackRequest{
			State: &state,
			Code:  &code,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("%s/app", h.config.FrontEndHost), http.StatusTemporaryRedirect)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	name := r.PostFormValue("name")
	birthday := r.PostFormValue("birthday")
	gender := r.PostFormValue("gender")
	sexualOrientation := r.PostFormValue("sexualOrientation")
	description := r.PostFormValue("description")

	birthdayDate, err := time.Parse(time.DateOnly, birthday)
	if err != nil {
		slog.Error("invalid birthday date", slog.String("error", err.Error()))
		http.Error(w, "invalid birthday date", http.StatusBadRequest)
		return
	}
	genderEnum := userServiceGrpc.Gender(userServiceGrpc.Gender_value[strings.ToUpper(gender)])
	sexualOrientationEnum := userServiceGrpc.SexualOrientation(userServiceGrpc.SexualOrientation_value[strings.ToUpper(sexualOrientation)])

	resp, err := h.UserServiceClient.CreateUser(
		r.Context(),
		&userServiceGrpc.CreateUserRequest{
			Email:             &email,
			Name:              &name,
			Birthday:          timestamppb.New(birthdayDate),
			Gender:            &genderEnum,
			SexualOrientation: &sexualOrientationEnum,
			Description:       &description,
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := &CreateUserResponse{
		UserID: resp.GetUserID(),
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonUser)
}

func (h *Handler) StorePhotos(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["photos[]"]

	photos := make([]*userServiceGrpc.Photo, len(files))
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "No se pudo abrir el archivo", http.StatusBadRequest)
			return
		}
		defer file.Close()

		photo, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "No se pudo abrir el archivo", http.StatusBadRequest)
			return
		}
		photos[i] = &userServiceGrpc.Photo{
			Filename: &fileHeader.Filename,
			Data:     photo,
		}
	}
	if len(files) == 0 {
		http.Error(w, "no se enviaron archivos", http.StatusBadRequest)
		return
	}
	_, err = h.UserServiceClient.StoreUserPhotos(
		r.Context(),
		&userServiceGrpc.StoreUserPhotosRequest{
			UserID: &userID,
			Photos: photos,
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type Track struct {
	Name    string   `json:"name"`
	Album   Album    `json:"album"`
	Artists []string `json:"artists"`
}

type Album struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Cover string `json:"cover"`
}

type Artist struct {
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
	Image  string   `json:"image"`
}

type GetUserRecommendationResponse struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Age        int32    `json:"age"`
	Bio        string   `json:"bio"`
	Photos     []string `json:"photos"`
	TopTracks  []Track  `json:"top_tracks"`
	TopArtists []Artist `json:"top_artists"`
}

func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_ = params["user_id"]

	recommendedUsers := []string{"6e71c9a8-cb7a-4e1e-b2f0-d25439a7802b"}
	getUserRecommendationsResponse := make([]GetUserRecommendationResponse, len(recommendedUsers))
	for i, user := range recommendedUsers {
		userResponse, err := h.UserServiceClient.GetUser(
			r.Context(),
			&userServiceGrpc.GetUserRequest{
				UserID: &user,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		topTracks := make([]Track, len(userResponse.GetTopTracks()))
		for j, track := range userResponse.GetTopTracks() {
			topTracks[j] = Track{
				Name: track.GetName(),
				Album: Album{
					Name:  track.GetAlbum().GetName(),
					Type:  track.GetAlbum().GetType(),
					Cover: track.GetAlbum().GetCover(),
				},
				Artists: track.GetArtists(),
			}
		}

		topArtists := make([]Artist, len(userResponse.GetTopArtists()))
		for j, artist := range userResponse.GetTopArtists() {
			topArtists[j] = Artist{
				Name:   artist.GetName(),
				Genres: artist.GetGenres(),
				Image:  artist.GetImage(),
			}
		}

		getUserRecommendationsResponse[i] = GetUserRecommendationResponse{
			Id:         userResponse.GetUserID(),
			Name:       userResponse.GetName(),
			Bio:        userResponse.GetDescription(),
			Photos:     userResponse.GetPhotos(),
			TopTracks:  topTracks,
			TopArtists: topArtists,
			Age:        userResponse.GetAge(),
		}
	}

	jsonGetUserRecommendationsResponse, err := json.Marshal(getUserRecommendationsResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonGetUserRecommendationsResponse)
}
