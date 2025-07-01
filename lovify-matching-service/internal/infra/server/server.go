package server

import (
	"context"
	"fmt"
	"github.com/gorkagg10/lovify/lovify-matching-service/internal/domain/profile"
	"github.com/gorkagg10/lovify/lovify-matching-service/internal/domain/recommender"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	matchingServiceGrpc "github.com/gorkagg10/lovify/lovify-matching-service/grpc/matching-service"
)

type MatchingServer struct {
	matchingServiceGrpc.UnimplementedMatchingServiceServer
	UserProfileCollection *mongo.Collection
	MusicProviderDataCollection *mongo.Collection
}

func NewMatchingServer(userProfileCollection *mongo.Collection) *MatchingServer {
	return &MatchingServer{
		UserProfileCollection: userProfileCollection,
	}
}

func (m *MatchingServer) RecommendUsers(ctx context.Context, req *matchingServiceGrpc.RecommendUsersRequest) (*matchingServiceGrpc.RecommendUsersResponse, error) {
	userID := req.GetUserID()

	var requester profile.UserProfile
	err := s.db.Collection("users").FindOne(ctx, bson.M{"email": userID}).Decode(&requester)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %v", err)
	}
	if requester.MusicProviderInfo() == nil {
		return nil, fmt.Errorf("usuario sin gustos musicales")
	}
	requestVector := recommender.BuildGenreVector(requester.)

	return &matchingServiceGrpc.RecommendUsersResponse{}, nil
}
