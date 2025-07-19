package server

import (
	"context"
	"encoding/base64"
	"fmt"
	userServiceGrpc "github.com/gorkagg10/lovify/lovify-user-service/grpc/user-service"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/domain/oauth"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/domain/profile"
	"github.com/gorkagg10/lovify/lovify-user-service/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"os"
	"path/filepath"
	"time"
)

type UserServer struct {
	userServiceGrpc.UnimplementedUserServiceServer
	ProfileManager *profile.Manager
	OAuthService   *oauth.Service
	uploadsDir     string
}

func NewUserServer(profileManager *profile.Manager, oAuthService *oauth.Service, uploadsDir string) *UserServer {
	return &UserServer{
		ProfileManager: profileManager,
		OAuthService:   oAuthService,
		uploadsDir:     uploadsDir,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userServiceGrpc.CreateUserRequest) (*userServiceGrpc.CreateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	userProfileID, err := s.ProfileManager.CreateUserProfile(ctx, req)
	if err != nil {
		return nil, err
	}
	return &userServiceGrpc.CreateUserResponse{
		UserID: &userProfileID,
	}, nil
}

func (s *UserServer) MusicProviderLogin(_ context.Context, req *userServiceGrpc.MusicProviderLoginRequest) (*userServiceGrpc.MusicProviderLoginResponse, error) {
	url := s.OAuthService.RequestAuthorization(req.GetUserID())
	return &userServiceGrpc.MusicProviderLoginResponse{
		Url: &url,
	}, nil
}

func (s *UserServer) MusicProviderOAuthCallback(ctx context.Context, req *userServiceGrpc.MusicProviderOAuthCallbackRequest) (*emptypb.Empty, error) {
	token, err := s.OAuthService.Exchange(ctx, req.GetCode())
	if err != nil {
		return nil, err
	}
	err = s.ProfileManager.ConnectWithMusicProvider(ctx, req.GetState(), token)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserServer) StoreUserPhotos(_ context.Context, req *userServiceGrpc.StoreUserPhotosRequest) (*emptypb.Empty, error) {
	uploadDir := filepath.Join(s.uploadsDir, req.GetUserID())
	err := os.RemoveAll(uploadDir)
	if err != nil {
		return nil, err
	}
	for _, photo := range req.GetPhotos() {
		err = storePhoto(uploadDir, photo.GetFilename(), photo.GetData())
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

func storePhoto(uploadDir string, filename string, data []byte) error {
	destinationPath := filepath.Join(uploadDir, filename)
	if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(destinationPath, data, 0644); err != nil {
		return fmt.Errorf("error guardando archivo: %v", err)
	}

	return nil
}

func getTopTracks(tracks []profile.Track) []*userServiceGrpc.Track {
	topTracks := make([]*userServiceGrpc.Track, len(tracks))
	for i, track := range tracks {
		topTracks[i] = &userServiceGrpc.Track{
			Name:    util.ValueToPointer(track.Name()),
			Artists: track.Artists(),
			Album: &userServiceGrpc.Album{
				Name:  util.ValueToPointer(track.Album().Type()),
				Type:  util.ValueToPointer(track.Album().Type()),
				Cover: util.ValueToPointer(track.Album().Image().URL()),
			},
		}
	}

	return topTracks
}

func getTopArtists(artists []profile.Artist) []*userServiceGrpc.Artist {
	topArtists := make([]*userServiceGrpc.Artist, len(artists))
	for i, artist := range artists {
		topArtists[i] = &userServiceGrpc.Artist{
			Name:   util.ValueToPointer(artist.Name()),
			Genres: artist.Genres(),
			Image:  util.ValueToPointer(artist.Image().URL()),
		}
	}

	return topArtists
}

func (s *UserServer) getPhotos(userID string) ([]string, error) {
	userDirectory := filepath.Join(s.uploadsDir, userID)
	photoEntries, err := os.ReadDir(userDirectory)
	if err != nil {
		return nil, err
	}
	photos := make([]string, len(photoEntries))

	for i, photoEntry := range photoEntries {
		photoPath := filepath.Join(userDirectory, photoEntry.Name())
		photo, err := readPhoto(photoPath)
		if err != nil {
			return nil, err
		}

		photos[i] = photo
	}

	return photos, nil
}

func readPhoto(photoPath string) (string, error) {
	file, err := os.Open(photoPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	photo, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(photo), nil
}

func calculateAge(birthday time.Time) int32 {
	today := time.Now()
	age := today.Year() - birthday.Year()

	if today.Month() < birthday.Month() || (today.Month() == birthday.Month() && today.Day() < birthday.Day()) {
		age--
	}

	return int32(age)
}

func (s *UserServer) GetUser(ctx context.Context, req *userServiceGrpc.GetUserRequest) (*userServiceGrpc.GetUserResponse, error) {
	userProfile, err := s.ProfileManager.GetProfile(ctx, req.GetUserID())
	if err != nil {
		return nil, err
	}
	gender := userServiceGrpc.Gender(userServiceGrpc.Gender_value[userProfile.Gender()])
	sexualOrientation := userServiceGrpc.SexualOrientation(userServiceGrpc.SexualOrientation_value[userProfile.SexualOrientation()])
	photos, err := s.getPhotos(req.GetUserID())
	if err != nil {
		return nil, err
	}
	age := calculateAge(userProfile.Birthday())
	return &userServiceGrpc.GetUserResponse{
		UserID:            req.UserID,
		Name:              util.ValueToPointer(userProfile.Name()),
		Description:       util.ValueToPointer(userProfile.Description()),
		Gender:            &gender,
		SexualOrientation: &sexualOrientation,
		Photos:            photos,
		TopTracks:         getTopTracks(userProfile.MusicProviderInfo().TopTracks()),
		TopArtists:        getTopArtists(userProfile.MusicProviderInfo().TopArtists()),
		Age:               util.ValueToPointer(age),
	}, nil
}
