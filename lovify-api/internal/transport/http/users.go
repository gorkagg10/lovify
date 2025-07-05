package http

import (
	"encoding/json"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	http.Redirect(w, r, "http://localhost:3000/app", http.StatusTemporaryRedirect)
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
