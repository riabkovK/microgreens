package domain

import "errors"

var (
	ErrUserNotFound           = errors.New("user doesn't exists")
	ErrUserAlreadyExists      = errors.New("user with such email already exists")
	ErrRefreshTokenHasExpired = errors.New("refresh token has been expired")
)
