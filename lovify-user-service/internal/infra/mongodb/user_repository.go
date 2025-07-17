package mongodb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorkagg10/lovify/lovify-user-service/internal/domain/oauth"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/domain/profile"
)

type UserRepository struct {
	UserProfileCollection         *mongo.Collection
	MusicProviderTokensCollection *mongo.Collection
	MusicProviderDataCollection   *mongo.Collection
}

func NewUserRepository(
	userProfileCollection *mongo.Collection,
	musicProviderTokensCollection *mongo.Collection,
	musicProviderDataCollection *mongo.Collection) *UserRepository {
	return &UserRepository{
		UserProfileCollection:         userProfileCollection,
		MusicProviderTokensCollection: musicProviderTokensCollection,
		MusicProviderDataCollection:   musicProviderDataCollection,
	}
}

func (u *UserRepository) CreateUserProfile(ctx context.Context, profile *profile.UserProfile) (string, error) {
	userProfileID := uuid.New().String()
	userProfile := NewUserProfile(
		userProfileID,
		profile.Email(),
		profile.Name(),
		profile.Birthday().String(),
		profile.Gender(),
		profile.SexualOrientation(),
		profile.Description(),
		profile.ConnectedToMusicProvider(),
	)
	_, err := u.UserProfileCollection.InsertOne(ctx, userProfile)
	if err != nil {
		return "", err
	}
	return userProfileID, nil
}

func (u *UserRepository) GetUserProfile(ctx context.Context, userID string) (*profile.UserProfile, error) {
	filter := bson.D{{"_id", userID}}

	var userProfile UserProfile
	err := u.UserProfileCollection.FindOne(ctx, filter).Decode(&userProfile)
	if err != nil {
		return nil, err
	}

	musicProviderDataFilter := bson.D{{"user_id", userID}}

	var musicProviderData MusicProviderData
	err = u.MusicProviderDataCollection.FindOne(ctx, musicProviderDataFilter).Decode(&musicProviderData)
	if err != nil {
		return nil, err
	}

	birthday, err := time.Parse("2006-01-02 15:04:05 -0700 MST", userProfile.Birthday)
	if err != nil {
		return nil, err
	}

	topTracks := make([]profile.Track, 5)
	for i, track := range musicProviderData.TopTracks[:5] {
		image := profile.NewImage(
			track.Album.Cover.Url,
			track.Album.Cover.Height,
			track.Album.Cover.Width,
		)
		album := profile.NewAlbum(
			track.Album.Name,
			track.Album.Type,
			image,
		)
		topTracks[i] = *profile.NewTrack(
			track.Name,
			album,
			track.Artists,
		)
	}

	topArtists := make([]profile.Artist, 5)
	for i, artist := range musicProviderData.TopArtist[:5] {
		image := profile.NewImage(
			artist.Image.Url,
			artist.Image.Height,
			artist.Image.Width,
		)
		topArtists[i] = *profile.NewArtist(
			artist.Name,
			artist.Genres,
			image,
		)
	}

	profileMusicProviderData := profile.NewMusicProviderData(topTracks, topArtists)

	return profile.NewUserProfile(
		userProfile.Email,
		birthday,
		userProfile.Name,
		userProfile.Gender,
		userProfile.SexualOrientation,
		userProfile.Description,
		userProfile.MusicProviderConnected,
		profileMusicProviderData,
	), nil
}

func (u *UserRepository) ConnectWithMusicProvider(ctx context.Context, userID string) error {
	filter := bson.D{{"_id", userID}}
	update := bson.D{{"$set", bson.D{{"music_provider_connected", true}}}}
	_, err := u.UserProfileCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) StoreMusicProviderToken(ctx context.Context, userID string, accessToken *oauth.Token) error {
	musicProviderToken := NewMusicProviderToken(
		userID,
		accessToken.AccessToken(),
		accessToken.RefreshToken(),
		accessToken.ExpiresAt(),
	)
	_, err := u.MusicProviderTokensCollection.InsertOne(ctx, musicProviderToken)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) StoreMusicProviderData(ctx context.Context, userID string, musicProviderData *profile.MusicProviderData) error {
	topTracks := make([]Track, len(musicProviderData.TopTracks()))
	for i, track := range musicProviderData.TopTracks() {
		albumCover := NewImage(track.Album().Image().URL(), track.Album().Image().Height(), track.Album().Image().Width())
		album := NewAlbum(track.Album().Name(), track.Album().Type(), albumCover)
		topTracks[i] = *NewTrack(track.Name(), album, track.Artists())
	}

	topArtists := make([]Artist, len(musicProviderData.TopArtists()))
	for i, artist := range musicProviderData.TopArtists() {
		image := NewImage(artist.Image().URL(), artist.Image().Height(), artist.Image().Width())
		topArtists[i] = *NewArtist(artist.Name(), artist.Genres(), image)
	}

	data := NewMusicProviderData(userID, topTracks, topArtists)
	_, err := u.MusicProviderDataCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
