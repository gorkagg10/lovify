package server

import (
	"context"
	messagingServiceGrpc "github.com/gorkagg10/lovify/lovify-messaging-service/grpc/messaging-service"
	"github.com/gorkagg10/lovify/lovify-messaging-service/internal/infra/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	filter := bson.M{"id": req.GetMatchID()}

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
		Content:    req.GetContent(),
		Read:       false,
	}

	_, err = m.MessageCollection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	if !match.ConversationStarted {
		update := bson.M{
			"$set": bson.M{"conversation_started": true},
		}

		_, err = m.MatchCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

func (m *MessagingServer) ListMessages(ctx context.Context, req *messagingServiceGrpc.ListMessagesRequest) (*messagingServiceGrpc.ListMessagesResponse, error) {
	messageFilter := bson.M{
		"match_id": req.GetMatchID(),
	}
	cursor, err := m.MessageCollection.Find(ctx, messageFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*messagingServiceGrpc.Message
	for cursor.Next(ctx) {
		var message mongodb.Message
		if err = cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages,
			&messagingServiceGrpc.Message{
				MatchID:    &message.MatchID,
				FromUserID: &message.FromUserID,
				ToUserID:   &message.ToUserID,
				Content:    &message.Content,
				SendAt:     timestamppb.New(message.SentAt),
				Read:       &message.Read,
			})
	}
	return &messagingServiceGrpc.ListMessagesResponse{
		Messages: messages,
	}, nil
}

func (m *MessagingServer) ListConversations(ctx context.Context, req *messagingServiceGrpc.ListConversationsRequest) (*messagingServiceGrpc.ListConversationsResponse, error) {
	matchFilter := bson.M{
		"$or": []bson.M{
			{"user_1_id": req.GetUserID()},
			{"user_2_id": req.GetUserID()},
		},
		"conversation_started": true,
	}
	cursor, err := m.MatchCollection.Find(ctx, matchFilter, options.Find().
		SetSort(bson.D{{Key: "matchedAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []*messagingServiceGrpc.Conversation
	for cursor.Next(ctx) {
		var match mongodb.Match
		if err = cursor.Decode(&match); err != nil {
			return nil, err
		}
		var message mongodb.Message
		opts := options.FindOne().
			SetSort(bson.D{{Key: "sent_at", Value: -1}}) // orden descendente

		err = m.MessageCollection.FindOne(ctx, bson.M{"match_id": match.ID}, opts).Decode(&message)
		if err != nil {
			return nil, err
		}
		otherID := match.User2ID
		if match.User2ID == req.GetUserID() {
			otherID = match.User1ID
		}
		conversations = append(conversations, &messagingServiceGrpc.Conversation{
			MatchID:     &match.ID,
			UserID:      &otherID,
			MatchedAt:   timestamppb.New(match.MatchedAt),
			LastMessage: &message.Content,
		})
	}
	return &messagingServiceGrpc.ListConversationsResponse{
		Conversations: conversations,
	}, nil
}
