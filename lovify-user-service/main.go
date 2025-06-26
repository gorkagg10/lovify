package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type SpotifySession struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type TopTracks struct {
	Tracks []Tracks `json:"items"`
}

type Tracks struct {
	Name    string    `json:"name"`
	Album   Album     `json:"album"`
	Artists []Artists `json:"artists"`
}

type Album struct {
	Name string `json:"name"`
}

type Artists struct {
	Name string `json:"name"`
}

type TopTracksResponse struct {
	Tracks []TracksResponse
}

type TracksResponse struct {
	Name    string   `json:"name"`
	Album   string   `json:"album"`
	Artists []string `json:"artists"`
}

type OAuthService struct {
	config *oauth2.Config
	state  string
}

func generateRandomState() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Intn(1000))
}

func NewOAuthService() *OAuthService {
	return &OAuthService{
		config: &oauth2.Config{
			ClientID:     "f4ed25e807ab4b74b981cd606a75699b",
			ClientSecret: "4b8515bf00ed4f67bbcd9a77d7486bdb",
			Endpoint:     spotify.Endpoint,
			RedirectURL:  "http://127.0.0.1:8082/callback/spotify",
			Scopes:       []string{"user-read-email", "user-read-recently-played", "user-top-read"},
		},
	}
}

func (o *OAuthService) HandleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	state := base64.URLEncoding.EncodeToString([]byte(userID))
	url := o.config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (o *OAuthService) HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	userIDBytes, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		http.Error(w, "Failed to decode state: "+err.Error(), http.StatusBadRequest)
	}
	userID := string(userIDBytes)
	fmt.Println("userID:", userID)

	code := r.URL.Query().Get("code")
	if state == "" {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := o.config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to exchange token:", err.Error())
		return
	}
	session := &SpotifySession{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}
	response, err := json.Marshal(session)
	if err != nil {
		http.Error(w, "Failed to marshal JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func GetTopTracks(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks", nil)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.Header.Set("time_range", "short_term")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read and print response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println(string(body))

	var topTracks TopTracks
	if err = json.Unmarshal(body, &topTracks); err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	response, err := formatResponse(&topTracks)
	if err != nil {
		fmt.Println("Error marshalling response body:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(response)
}

func formatResponse(tracks *TopTracks) ([]byte, error) {
	topTracks := make([]TracksResponse, len(tracks.Tracks))
	for i, track := range tracks.Tracks {
		topTracks[i] = TracksResponse{
			Name:  track.Name,
			Album: track.Album.Name,
		}
		artists := make([]string, len(tracks.Tracks[i].Artists))
		for j, artist := range tracks.Tracks[i].Artists {
			artists[j] = artist.Name
		}
		topTracks[i].Artists = artists
	}
	return json.Marshal(topTracks)
}

func main() {
	router := mux.NewRouter()

	oauthService := NewOAuthService()

	router.HandleFunc("/users/{user_id}/login/spotify", oauthService.HandleSpotifyLogin)
	router.HandleFunc("/callback/spotify", oauthService.HandleCallback)
	router.HandleFunc("/spotify/me", GetTopTracks)

	server := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					fmt.Println("Fetching top tracks")
				}
			}
		}(context.Background())
	}()

	fmt.Println("Server is running on http://localhost:8082")
	log.Fatal(server.ListenAndServe())
}
