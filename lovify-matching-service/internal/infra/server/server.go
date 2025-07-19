package server

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/gorkagg10/lovify/lovify-matching-service/internal/domain/recommender"
	"github.com/gorkagg10/lovify/lovify-matching-service/internal/infra/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sort"
	"time"

	matchingServiceGrpc "github.com/gorkagg10/lovify/lovify-matching-service/grpc/matching-service"
)

const (
	minScore = 0.15
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

func (m *MatchingServer) getUsers(ctx context.Context) ([]recommender.User, error) {
	cursor, err := m.UserProfileCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var profiles []recommender.User
	for cursor.Next(ctx) {
		var profile mongodb.UserProfile
		if err = cursor.Decode(&profile); err != nil {
			continue
		}

		var musicData mongodb.MusicProviderData
		err = m.MusicProviderDataCollection.FindOne(ctx, bson.M{"user_id": profile.ID}).Decode(&musicData)

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

		topTracks := make([]recommender.Track, len(musicData.TopTracks))
		for i, topTrack := range musicData.TopTracks {
			album := recommender.Album{
				Name:      topTrack.Album.Name,
				AlbumType: topTrack.Album.Type,
				Image: &recommender.Image{
					Url:    topTrack.Album.Cover.Url,
					Height: topTrack.Album.Cover.Height,
					Width:  topTrack.Album.Cover.Width,
				},
			}
			topTracks[i] = recommender.Track{
				Name:    topTrack.Name,
				Album:   &album,
				Artists: topTrack.Artists,
			}
		}

		birthday, err := time.Parse("2006-01-02 15:04:05 -0700 MST", profile.Birthday)
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, recommender.User{
			ID:                       profile.ID,
			Email:                    profile.Email,
			Name:                     profile.Name,
			Birthday:                 birthday,
			Gender:                   profile.Gender,
			SexualOrientation:        profile.SexualOrientation,
			Description:              profile.Description,
			ConnectedToMusicProvider: profile.MusicProviderConnected,
			MusicProviderInfo: &recommender.MusicProviderData{
				TopArtists: topArtists,
				TopTracks:  topTracks,
			},
		})
	}

	return profiles, nil
}

// BuildPreferences creates for each user an ordered list of compatible candidates.
func (m *MatchingServer) BuildPreferences(ctx context.Context, users []recommender.User) map[string][]string {
	preferences := map[string][]string{}
	vectorCache := map[string]map[string]float64{}

	getVector := func(u *recommender.User) map[string]float64 {
		if v, ok := vectorCache[u.ID]; ok {
			return v
		}
		vectorCache[u.ID] = recommender.GenreVector(u, 5)
		return vectorCache[u.ID]
	}

	compatible := func(a, b recommender.User) bool {
		switch {
		case a.SexualOrientation == "HETEROSEXUAL" && b.SexualOrientation == "HETEROSEXUAL":
			return a.Gender != b.Gender
		case a.SexualOrientation == "HOMOSEXUAL" && b.SexualOrientation == "HOMOSEXUAL":
			return a.Gender == b.Gender
		default:
			return true // simplificación para "bi" o mixto
		}
	}

	liked := func(a, b recommender.User) bool {
		likeFilter := bson.M{
			"from_user_id": a.ID,
			"to_user_id":   b.ID,
		}
		err := m.LikeCollection.FindOne(ctx, likeFilter).Err()
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			return false
		}
		return true
	}

	for _, u := range users {
		var list []struct {
			id    string
			score float64
		}
		ua := getVector(&u)
		for _, v := range users {
			if u.ID == v.ID || !compatible(u, v) || liked(u, v) {
				continue
			}
			uv := getVector(&v)
			sim := recommender.CosSim(ua, uv)
			if sim >= minScore {
				list = append(list, struct {
					id    string
					score float64
				}{v.ID, sim})
			}
		}
		sort.Slice(list, func(i, j int) bool { return list[i].score > list[j].score })
		for _, it := range list {
			preferences[u.ID] = append(preferences[u.ID], it.id)
		}
	}
	return preferences
}

func (m *MatchingServer) Recommend(ctx context.Context, users []recommender.User, requester string) (string, float64) {
	prefs := m.BuildPreferences(ctx, users)
	if _, ok := prefs[requester]; !ok || len(prefs[requester]) == 0 {
		return "", 0
	}

	// GS simplificado (proposals solo del requester hasta conseguir match estable)
	current := map[string]string{}
	next := 0
	prop := requester
	for {
		if next >= len(prefs[prop]) {
			break
		}
		target := prefs[prop][next]
		next++
		if _, taken := current[target]; !taken {
			current[target] = prop
			break
		}
	}
	if receiver, ok := current[prefs[requester][next-1]]; ok && receiver == requester {
		// devuelve el score real
		vecReq := recommender.GenreVector(findUser(users, requester), 5)
		vecRec := recommender.GenreVector(findUser(users, prefs[requester][next-1]), 5)
		return prefs[requester][next-1], recommender.CosSim(vecReq, vecRec)
	}
	return "", 0
}

func findUser(users []recommender.User, id string) *recommender.User {
	for i := range users {
		if users[i].ID == id {
			return &users[i]
		}
	}
	return nil
}

func (m *MatchingServer) RecommendUser(ctx context.Context, req *matchingServiceGrpc.RecommendUserRequest) (*matchingServiceGrpc.RecommendUserResponse, error) {
	userID := req.GetUserID()
	users, err := m.getUsers(ctx)
	if err != nil {
		return nil, err
	}

	recommendedUserID, _ := m.Recommend(ctx, users, userID)

	return &matchingServiceGrpc.RecommendUserResponse{
		RecommendedUserID: &recommendedUserID,
	}, nil
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
