package mongodb

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorkagg10/lovify-user-service/internal/domain/oauth"
	"github.com/gorkagg10/lovify-user-service/internal/domain/profile"
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
