package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorkagg10/lovify/lovify-matching-service/internal/domain/profile"
)

type UserRepository struct {
	UserProfileCollection       *mongo.Collection
	MusicProviderDataCollection *mongo.Collection
}

func NewUserRepository(
	userProfileCollection *mongo.Collection,
	musicProviderDataCollection *mongo.Collection) *UserRepository {
	return &UserRepository{
		UserProfileCollection:       userProfileCollection,
		MusicProviderDataCollection: musicProviderDataCollection,
	}
}

func (u *UserRepository) GetUserProfile(ctx context.Context, userID string) (*profile.UserProfile, error) {
	var user
	err := u.UserProfileCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&requester)
	if err != nil {
		return "", err
	}
	return userProfileID, nil
}
