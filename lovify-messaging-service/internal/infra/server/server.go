package server

import (
	"context"
	messagingServiceGrpc "github.com/gorkagg10/lovify/lovify-messaging-service/grpc/messaging-service"
	"github.com/gorkagg10/lovify/lovify-messaging-service/internal/infra/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type MessagingServer struct {
	messagingServiceGrpc.UnimplementedMessagingServiceServer
	MatchCollection   *mongo.Collection
	MessageCollection *mongo.Collection
}

func NewMessagingServer(
	matchCollection *mongo.Collection,
	messageCollection *mongo.Collection,
) *MessagingServer {
	return &MessagingServer{
		MatchCollection:   matchCollection,
		MessageCollection: messageCollection,
	}
}

func (m *MessagingServer) SendMessage(ctx context.Context, req *messagingServiceGrpc.SendMessageRequest) (*emptypb.Empty, error) {
	var match mongodb.Match

	filter := bson.M{"matchId": req.GetMatchID()}

	err := m.MatchCollection.FindOne(ctx, filter).Decode(&match)
	if err != nil {
		return nil, err
	}

	toUserID := match.User2ID
	if toUserID == req.GetUserID() {
		toUserID = match.User1ID
	}

	message := mongodb.Message{
		MatchID:    req.GetMatchID(),
		FromUserID: req.GetUserID(),
		ToUserID:   toUserID,
		SentAt:     time.Now(),
	}

	_, err = m.MessageCollection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	if !match.ConversationStarted {
		update := bson.M{
			"$set": bson.M{"conversationStarted": true},
		}

		_, err = m.MatchCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}
