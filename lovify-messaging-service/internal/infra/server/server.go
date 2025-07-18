package server

import (
	messagingServiceGrpc "github.com/gorkagg10/lovify/lovify-messaging-service/grpc/messaging-service"
	"go.mongodb.org/mongo-driver/mongo"
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
