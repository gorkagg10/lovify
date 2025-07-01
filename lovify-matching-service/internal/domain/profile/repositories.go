package profile

import (
	"context"
)

type UserRepository interface {
	GetUserProfile(ctx context.Context, userID string) (*UserProfile, error)
}
