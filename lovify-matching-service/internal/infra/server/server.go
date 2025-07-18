package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorkagg10/lovify/lovify-matching-service/internal/domain/recommender"
	"github.com/gorkagg10/lovify/lovify-matching-service/internal/infra/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	matchingServiceGrpc "github.com/gorkagg10/lovify/lovify-matching-service/grpc/matching-service"
)

const (
	minSimilarity = 0
)

type MatchingServer struct {
	matchingServiceGrpc.UnimplementedMatchingServiceServer
	UserProfileCollection       *mongo.Collection
	MusicProviderDataCollection *mongo.Collection
	LikeCollection              *mongo.Collection
	MatchCollection             *mongo.Collection
}

func NewMatchingServer(
	userProfileCollection *mongo.Collection,
	musicProviderDataCollection *mongo.Collection,
	likeCollection *mongo.Collection,
	matchCollection *mongo.Collection,
) *MatchingServer {
	return &MatchingServer{
		UserProfileCollection:       userProfileCollection,
		MusicProviderDataCollection: musicProviderDataCollection,
		LikeCollection:              likeCollection,
		MatchCollection:             matchCollection,
	}
}

func (m *MatchingServer) RecommendUsers(ctx context.Context, req *matchingServiceGrpc.RecommendUsersRequest) (*matchingServiceGrpc.RecommendUsersResponse, error) {
	userID := req.GetUserID()

	var requester mongodb.UserProfile
	err := m.UserProfileCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&requester)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %v", err)
	}

	var musicData mongodb.MusicProviderData
	err = m.MusicProviderDataCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&musicData)

	topArtists := make([]recommender.Artist, len(musicData.TopArtist))
	for i, topArtist := range musicData.TopArtist {
		topArtists[i] = recommender.Artist{
			Name:   topArtist.Name,
			Genres: topArtist.Genres,
			Image: &recommender.Image{
				Url:    topArtist.Image.Url,
				Height: topArtist.Image.Height,
				Width:  topArtist.Image.Width,
			},
		}
	}
	requestVector := recommender.BuildGenreVector(topArtists)
	// Cargar usuarios candidatos
	cursor, err := m.UserProfileCollection.Find(ctx, bson.M{"email": bson.M{"$ne": userID}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var recommendedUserIDs []string
	for cursor.Next(ctx) {
		var candidate mongodb.UserProfile
		if err = cursor.Decode(&candidate); err != nil {
			continue
		}

		var candidateMusicData mongodb.MusicProviderData
		err = m.MusicProviderDataCollection.FindOne(ctx, bson.M{"user_id": candidate.ID}).Decode(&candidateMusicData)
		if err != nil || !compatible(requester, candidate) {
			continue
		}

		candidateTopArtists := make([]recommender.Artist, len(candidateMusicData.TopArtist))
		for i, topArtist := range candidateMusicData.TopArtist {
			candidateTopArtists[i] = recommender.Artist{
				Name:   topArtist.Name,
				Genres: topArtist.Genres,
				Image: &recommender.Image{
					Url:    topArtist.Image.Url,
					Height: topArtist.Image.Height,
					Width:  topArtist.Image.Width,
				},
			}
		}

		candidateVector := recommender.BuildGenreVector(candidateTopArtists)
		score := recommender.CosSim(requestVector, candidateVector)
		if score >= minSimilarity {
			recommendedUserIDs = append(recommendedUserIDs, candidate.ID)
		}
	}

	return &matchingServiceGrpc.RecommendUsersResponse{
		RecommendedUsersIDs: recommendedUserIDs,
	}, nil
}

func compatible(userA, userB mongodb.UserProfile) bool {
	if userA.SexualOrientation == "HETEROSEXUAL" && userB.SexualOrientation == "HETEROSEXUAL" {
		return userA.Gender != userB.Gender
	}
	if userA.SexualOrientation == "HOMOSEXUAL" && userB.SexualOrientation == "HOMOSEXUAL" {
		return userA.Gender == userB.Gender
	}
	return true
}

func (m *MatchingServer) HandleLike(ctx context.Context, req *matchingServiceGrpc.HandleLikeRequest) (*emptypb.Empty, error) {
	newLike := mongodb.Like{
		FromUserId: req.GetFromUserID(),
		ToUserId:   req.GetToUserID(),
		Type:       req.GetType().String(),
		CreatedAt:  time.Now(),
	}
	_, err := m.LikeCollection.InsertOne(ctx, newLike)
	if err != nil {
		return nil, err
	}
	if req.GetType() == matchingServiceGrpc.Like_LIKE {
		matchFilter := bson.M{
			"from_user_id": req.GetToUserID(),
			"to_user_id":   req.GetFromUserID(),
			"type":         matchingServiceGrpc.Like_LIKE.String(),
		}
		var reverse mongodb.Like
		err = m.LikeCollection.FindOne(ctx, matchFilter).Decode(&reverse)
		if err == nil {
			matchID := uuid.New().String()
			match := mongodb.Match{
				ID:                  matchID,
				User1ID:             req.GetFromUserID(),
				User2ID:             req.GetToUserID(),
				MatchedAt:           time.Now(),
				ConversationStarted: false,
			}
			_, err = m.MatchCollection.InsertOne(ctx, match)
			if err != nil {
				return nil, err
			}
		}
	}
	return &emptypb.Empty{}, nil
}

func (m *MatchingServer) GetMatches(ctx context.Context, req *matchingServiceGrpc.GetMatchesRequest) (*matchingServiceGrpc.GetMatchesResponse, error) {
	matchFilter := bson.M{
		"$or": []bson.M{
			{"user_1_id": req.GetUserID()},
			{"user_2_id": req.GetUserID()},
		},
		"conversation_started": false,
	}
	cursor, err := m.MatchCollection.Find(ctx, matchFilter, options.Find().
		SetSort(bson.D{{Key: "matchedAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []*matchingServiceGrpc.Match
	for cursor.Next(ctx) {
		var match mongodb.Match
		if err = cursor.Decode(&match); err != nil {
			return nil, err
		}
		otherID := match.User2ID
		if match.User2ID == req.GetUserID() {
			otherID = match.User1ID
		}
		matches = append(matches, &matchingServiceGrpc.Match{
			MatchID:   &match.ID,
			UserID:    &otherID,
			MatchedAt: timestamppb.New(match.MatchedAt),
		})
	}
	return &matchingServiceGrpc.GetMatchesResponse{
		Matches: matches,
	}, nil
}
