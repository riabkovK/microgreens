package domain

import "time"

type Session struct {
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RefreshSessionCheckResponse struct {
	Id           int       `json:"id" db:"id"`
	UserId       int       `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
}
