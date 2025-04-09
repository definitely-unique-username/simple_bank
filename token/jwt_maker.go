package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const MIN_SECRET_KEY_SIZE = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) > MIN_SECRET_KEY_SIZE {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", MIN_SECRET_KEY_SIZE)
	}

	return &JWTMaker{secretKey}, nil
}

func (m *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(m.secretKey))
}

func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrInvlidToken
		}

		return []byte(m.secretKey), nil
	})

	if err != nil {
		// don't know how to check exact error
		// return nil, ErrInvalidToken

		return nil, ErrExpiredToken
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvlidToken
	}

	return payload, nil
}
