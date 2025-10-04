package model

import (
	"encoding/json"
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

type GroupId uuid.UUID
type Group struct {
	ID          GroupId   `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   UserId    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GroupDetails struct {
	Group   Group
	Members []*Gift `json:"membersGifts"`
}

type ModelId interface {
	UserId | SessionID | GiftId | GroupId
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
	GroupId     GroupId    `json:"groupId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Link        string     `json:"link"`
	Price       float64    `json:"price,omitempty"`
	CreatedBy   UserId     `json:"createdBy"`
	ClaimedBy   *UserId    `json:"claimedBy,omitempty"`
	ClaimedAt   *time.Time `json:"claimedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// Generic UUID marshaling helpers
func marshalUUID[T ModelId](id T) ([]byte, error) {
	return json.Marshal(uuid.UUID(id).String())
}

func unmarshalUUID[T ModelId](data []byte) (T, error) {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		var zero T
		return zero, err
	}
	parsed, err := uuid.Parse(str)
	if err != nil {
		var zero T
		return zero, err
	}
	return T(parsed), nil
}

func (u UserId) MarshalJSON() ([]byte, error) {
	return marshalUUID(u)
}

func (u *UserId) UnmarshalJSON(data []byte) error {
	parsed, err := unmarshalUUID[UserId](data)
	if err != nil {
		return err
	}
	*u = parsed
	return nil
}

func (g GroupId) MarshalJSON() ([]byte, error) {
	return marshalUUID(g)
}

func (g *GroupId) UnmarshalJSON(data []byte) error {
	parsed, err := unmarshalUUID[GroupId](data)
	if err != nil {
		return err
	}
	*g = parsed
	return nil
}

func (g GiftId) MarshalJSON() ([]byte, error) {
	return marshalUUID(g)
}

func (g *GiftId) UnmarshalJSON(data []byte) error {
	parsed, err := unmarshalUUID[GiftId](data)
	if err != nil {
		return err
	}
	*g = parsed
	return nil
}

func (s SessionID) MarshalJSON() ([]byte, error) {
	return marshalUUID(s)
}

func (s *SessionID) UnmarshalJSON(data []byte) error {
	parsed, err := unmarshalUUID[SessionID](data)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}
