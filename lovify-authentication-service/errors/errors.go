package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrUserAlreadyExistsMsg              = "USER_ALREADY_EXISTS"
	ErrUserNotFoundMsg                   = "USER_NOT_FOUND"
	ErrDatabaseQueryFailedMsg            = "DATABASE_QUERY_FAILED"
	ErrHashedPasswordGenerationFailedMsg = "HASHED_PASSWORD_GENERATION_FAILED"
	ErrIncorrectPasswordMsg              = "INCORRECT_PASSWORD"
)

var (
	StatusUserAlreadyExists              = status.New(codes.InvalidArgument, ErrUserAlreadyExistsMsg)
	StatusDatabaseQueryFailed            = status.New(codes.Internal, ErrDatabaseQueryFailedMsg)
	StatusHashedPasswordGenerationFailed = status.New(codes.Internal, ErrHashedPasswordGenerationFailedMsg)
	StatusUserNotFound                   = status.New(codes.NotFound, ErrUserNotFoundMsg)
	StatusIncorrectPassword              = status.New(codes.InvalidArgument, ErrIncorrectPasswordMsg)

	ErrUserAlreadyExists              = StatusUserAlreadyExists.Err()
	ErrDatabaseQueryFailed            = StatusDatabaseQueryFailed.Err()
	ErrHashedPasswordGenerationFailed = StatusHashedPasswordGenerationFailed.Err()
	ErrUserNotFound                   = StatusUserNotFound.Err()
	ErrIncorrectPassword              = StatusIncorrectPassword.Err()
)
