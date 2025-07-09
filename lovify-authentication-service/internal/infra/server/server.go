package server

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	authServiceGrpc "github.com/gorkagg10/lovify/lovify-authentication-service/grpc/auth-service"
	"github.com/gorkagg10/lovify/lovify-authentication-service/internal/domain/login"
	"github.com/gorkagg10/lovify/lovify-authentication-service/util"
)

type AuthServer struct {
	authServiceGrpc.UnimplementedAuthServiceServer
	authenticationService *login.Authorization
}

func NewAuthServer(authenticationService *login.Authorization) *AuthServer {
	return &AuthServer{
		authenticationService: authenticationService,
	}
}

func (s *AuthServer) RegisterUser(ctx context.Context, req *authServiceGrpc.RegisterRequest) (*emptypb.Empty, error) {
	err := s.authenticationService.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authServiceGrpc.LoginRequest) (*authServiceGrpc.LoginResponse, error) {
	user, err := s.authenticationService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &authServiceGrpc.LoginResponse{
		SessionToken: &authServiceGrpc.Token{
			Token:          util.ValueToPointer(user.SessionToken().Token()),
			ExpirationDate: timestamppb.New(user.SessionToken().ExpirationDate()),
		},
		CsrfToken: &authServiceGrpc.Token{
			Token:          util.ValueToPointer(user.CSRFToken().Token()),
			ExpirationDate: timestamppb.New(user.CSRFToken().ExpirationDate()),
		},
	}, nil
}

func (s *AuthServer) Authorize(ctx context.Context, req *authServiceGrpc.AuthorizationRequest) (*emptypb.Empty, error) {
	err := s.authenticationService.AuthorizeUser(ctx, req.GetEmail(), req.GetSessionToken(), req.GetCsrfToken())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
