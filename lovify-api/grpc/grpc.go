package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient[T any](target string, clientFn func(connInterface grpc.ClientConnInterface) T) (T, error) {
	conn, err := grpc.NewClient(target, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}...)
	if err != nil {
		var empty T
		return empty, err
	}
	return clientFn(conn), nil
}
