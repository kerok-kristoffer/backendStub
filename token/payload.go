package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var ErrInvalidToken = errors.New("token is invalid")
var ErrExpiredToken = errors.New("token has expired")

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expired_at"`
}

func (p Payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId := uuid.New()

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	return payload, nil
}
