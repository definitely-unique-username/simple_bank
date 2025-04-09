package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	issuedAt := time.Now()

	return &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: issuedAt.Add(duration),
	}, nil
}

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvlidToken  = errors.New("token is invalid")
)

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (p *Payload) GetIssuer() (string, error) {
	return "", nil
}

func (p *Payload) GetSubject() (string, error) {
	return "", nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	var claims []string

	return claims, nil
}
