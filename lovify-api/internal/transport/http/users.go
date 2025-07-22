package http

import (
	"encoding/json"
	"fmt"
	matchingServiceGrpc "github.com/gorkagg10/lovify/lovify-matching-service/grpc/matching-service"
	messagingServiceGrpc "github.com/gorkagg10/lovify/lovify-messaging-service/grpc/messaging-service"
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

type GetUserResponse struct {
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
	userID := params["user_id"]

	recommendedUser, err := h.MatchingServiceClient.RecommendUser(
		r.Context(),
		&matchingServiceGrpc.RecommendUserRequest{
			UserID: &userID,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if recommendedUser.GetRecommendedUserID() == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	userResponse, err := h.UserServiceClient.GetUser(
		r.Context(),
		&userServiceGrpc.GetUserRequest{
			UserID: recommendedUser.RecommendedUserID,
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

	getRecommendationResponse := GetUserResponse{
		Id:         userResponse.GetUserID(),
		Name:       userResponse.GetName(),
		Bio:        userResponse.GetDescription(),
		Photos:     userResponse.GetPhotos(),
		TopTracks:  topTracks,
		TopArtists: topArtists,
		Age:        userResponse.GetAge(),
	}

	jsonGetUserRecommendationsResponse, err := json.Marshal(getRecommendationResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonGetUserRecommendationsResponse)
}

type HandleLikeRequest struct {
	Type string `json:"type"`
}

func (h *Handler) HandleLike(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fromID := params["user_id"]
	toID := params["to_id"]

	var likeRequest HandleLikeRequest
	err := json.NewDecoder(r.Body).Decode(&likeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	likeType := matchingServiceGrpc.Like(matchingServiceGrpc.Like_value[strings.ToUpper(likeRequest.Type)])

	_, err = h.MatchingServiceClient.HandleLike(
		r.Context(),
		&matchingServiceGrpc.HandleLikeRequest{
			FromUserID: &fromID,
			ToUserID:   &toID,
			Type:       &likeType,
		},
	)
	w.WriteHeader(http.StatusOK)
}

type GetMatchesResponse struct {
	MatchID    string `json:"match_id"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	MatchedAt  string `json:"matched_at"`
	FirstImage string `json:"first_image"`
}

func (h *Handler) GetMatches(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	matchesInfo, err := h.MatchingServiceClient.GetMatches(
		r.Context(),
		&matchingServiceGrpc.GetMatchesRequest{
			UserID: &userID,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	matches := make([]GetMatchesResponse, len(matchesInfo.Matches))
	for i, match := range matchesInfo.Matches {
		userInfo, err := h.UserServiceClient.GetUser(
			r.Context(),
			&userServiceGrpc.GetUserRequest{
				UserID: match.UserID,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matches[i] = GetMatchesResponse{
			MatchID:    match.GetMatchID(),
			UserID:     match.GetUserID(),
			MatchedAt:  match.GetMatchedAt().AsTime().Format(time.RFC3339),
			Name:       userInfo.GetName(),
			FirstImage: userInfo.Photos[0],
		}
	}

	jsonGetMatchesResponse, err := json.Marshal(matches)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonGetMatchesResponse)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	userResponse, err := h.UserServiceClient.GetUser(
		r.Context(),
		&userServiceGrpc.GetUserRequest{
			UserID: &userID,
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

	user := GetUserResponse{
		Id:         userResponse.GetUserID(),
		Name:       userResponse.GetName(),
		Bio:        userResponse.GetDescription(),
		Photos:     userResponse.GetPhotos(),
		TopTracks:  topTracks,
		TopArtists: topArtists,
		Age:        userResponse.GetAge(),
	}

	getUserResponse, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(getUserResponse)
}

type GetConversationsResponse struct {
	MatchID    string `json:"match_id"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	MatchedAt  string `json:"matched_at"`
	FirstImage string `json:"first_image"`
}

func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	conversationsInfo, err := h.MessagingServiceClient.ListConversations(
		r.Context(),
		&messagingServiceGrpc.ListConversationsRequest{
			UserID: &userID,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conversations := make([]GetConversationsResponse, len(conversationsInfo.Conversations))
	for i, conversation := range conversationsInfo.Conversations {
		userInfo, err := h.UserServiceClient.GetUser(
			r.Context(),
			&userServiceGrpc.GetUserRequest{
				UserID: conversation.UserID,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		conversations[i] = GetConversationsResponse{
			MatchID:    conversation.GetMatchID(),
			UserID:     conversation.GetUserID(),
			MatchedAt:  conversation.GetMatchedAt().AsTime().Format(time.RFC3339),
			Name:       userInfo.GetName(),
			FirstImage: userInfo.Photos[0],
		}
	}

	jsonGetConversationsResponse, err := json.Marshal(conversations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonGetConversationsResponse)
}
