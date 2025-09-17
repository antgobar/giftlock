package model

import (
	"time"

	"github.com/google/uuid"
)

type UserId uuid.UUID
type User struct {
	ID             UserId    `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
}

type SessionID uuid.UUID
type SessionToken string
type Session struct {
	ID        SessionID    `json:"id"`
	UserId    UserId       `json:"userId"`
	Token     SessionToken `json:"token"`
	CreatedAt time.Time    `json:"createdAt"`
	ExpiresAt time.Time    `json:"expiresAt"`
}
