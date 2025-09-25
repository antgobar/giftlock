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

type ModelId interface {
	UserId | SessionID | GiftId
}

func IdFromString[T ModelId](id string) (T, error) {
	u, err := uuid.Parse(id)
	return T(u), err
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

type GiftId uuid.UUID
type Gift struct {
	ID          GiftId     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Link        string     `json:"link"`
	Price       float64    `json:"price,omitempty"`
	CreatedBy   UserId     `json:"createdBy"`
	ClaimedBy   *UserId    `json:"claimedBy,omitempty"`
	ClaimedAt   *time.Time `json:"claimedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}
