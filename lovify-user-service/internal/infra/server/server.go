package server

import (
	"context"
	userServiceGrpc "github.com/gorkagg10/lovify/lovify-user-service/grpc/user-service"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/domain/oauth"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/domain/profile"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	userServiceGrpc.UnimplementedUserServiceServer
	ProfileManager *profile.Manager
	OAuthService   *oauth.Service
}

func NewUserServer(profileManager *profile.Manager, oAuthService *oauth.Service) *UserServer {
	return &UserServer{
		ProfileManager: profileManager,
		OAuthService:   oAuthService,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userServiceGrpc.CreateUserRequest) (*userServiceGrpc.CreateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	userProfileID, err := s.ProfileManager.CreateUserProfile(ctx, req)
	if err != nil {
		return nil, err
	}
	return &userServiceGrpc.CreateUserResponse{
		UserID: &userProfileID,
	}, nil
}

func (s *UserServer) MusicProviderLogin(_ context.Context, req *userServiceGrpc.MusicProviderLoginRequest) (*userServiceGrpc.MusicProviderLoginResponse, error) {
	url := s.OAuthService.RequestAuthorization(req.GetUserID())
	return &userServiceGrpc.MusicProviderLoginResponse{
		Url: &url,
	}, nil
}

func (s *UserServer) MusicProviderOAuthCallback(ctx context.Context, req *userServiceGrpc.MusicProviderOAuthCallbackRequest) (*emptypb.Empty, error) {
	token, err := s.OAuthService.Exchange(ctx, req.GetCode())
	if err != nil {
		return nil, err
	}
	err = s.ProfileManager.ConnectWithMusicProvider(ctx, req.GetState(), token)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
