package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	userService "github.com/gorkagg10/lovify/lovify-user-service/events"
	"github.com/nats-io/nats.go/jetstream"
	"log/slog"
	"time"
)

const (
	SessionToken = "session"
	CSRFToken    = "csrf"
)

type Authorization struct {
	userRepository     UserRepository
	securityRepository SecurityRepository
	tokenRepository    TokenRepository
}

func NewAuthorization(
	userRepository UserRepository,
	securityRepository SecurityRepository,
	tokenRepository TokenRepository) *Authorization {
	return &Authorization{
		userRepository:     userRepository,
		securityRepository: securityRepository,
		tokenRepository:    tokenRepository,
	}
}

func (a *Authorization) Register(ctx context.Context, email string, password string) error {
	exists, err := a.userRepository.EmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}
	hashedPassword, err := a.securityRepository.HashPassword(password)
	if err != nil {
		return err
	}
	user := NewUser(email, hashedPassword, false, "")
	if err = a.userRepository.CreateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (a *Authorization) Login(ctx context.Context, email string, password string) (*User, error) {
	user, err := a.userRepository.GetUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	validPassword, err := a.securityRepository.CheckPassword(user.HashedPassword(), password)
	if !validPassword {
		return nil, fmt.Errorf("invalid password")
	}
	if err != nil {
		return nil, fmt.Errorf("checking password: %w", err)
	}

	sessionToken, err := a.tokenRepository.GenerateToken(SessionToken, user.Email())
	if err != nil {
		return nil, fmt.Errorf("generating session token: %w", err)
	}
	user.setSessionToken(sessionToken)

	csrfToken, err := a.tokenRepository.GenerateToken(CSRFToken, user.Email())
	if err != nil {
		return nil, fmt.Errorf("generating session token: %w", err)
	}

	user.setCSRFToken(csrfToken)

	return user, nil
}

func isValidToken(userToken *Token) bool {
	if userToken.ExpirationDate().After(time.Now()) {
		return false
	}
	return true
}

func (a *Authorization) AuthorizeUser(ctx context.Context, email string, sessionToken string, csrfToken string) error {
	err := a.tokenRepository.ValidateToken(sessionToken)
	if err != nil {
		return fmt.Errorf("validating session token: %w", err)
	}
	err = a.tokenRepository.ValidateToken(csrfToken)
	if err != nil {
		return fmt.Errorf("validating csrf token: %w", err)
	}
	return nil
}

func (a *Authorization) ProcessProfileCreation(msg jetstream.Msg) {
	profile := new(userService.Profile)
	if err := json.Unmarshal(msg.Data(), profile); err != nil {
		if err = msg.Nak(); err != nil {
			slog.Error("not acknowledging message", err)
		}
		return
	}
	err := a.userRepository.ConnectProfile(profile.Email, profile.UserID)
	if err != nil {
		slog.Error("error connecting profile to user", err)
		if err = msg.Nak(); err != nil {
			slog.Error("not acknowledging message", err)
		}
		return
	}
	if err = msg.Ack(); err != nil {
		slog.Error("error acknowledging message", err)
	}
}
