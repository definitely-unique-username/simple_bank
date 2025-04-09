package token

import (
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type PasetoMaker struct {
	symmetricalKey paseto.V4SymmetricKey
	implicit       []byte
}

func NewPasetoMaker(symmetricalKey string) Maker {
	return &PasetoMaker{
		symmetricalKey: paseto.NewV4SymmetricKey(),
		implicit:       []byte("implicit nonce"),
	}
}

func (m *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}

	token := paseto.NewToken()

	token.SetString("id", payload.ID.String())
	token.SetString("username", payload.Username)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetNotBefore(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	return token.V4Encrypt(m.symmetricalKey, m.implicit), nil
}

func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	parsedToken, err := parser.ParseV4Local(m.symmetricalKey, token, m.implicit)

	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrExpiredToken
		}

		return nil, err
	}

	payload, err := extractPayload(parsedToken)

	if err != nil {
		return nil, ErrInvlidToken
	}

	return payload, nil
}

func extractPayload(t *paseto.Token) (*Payload, error) {
	id, err := t.GetString("id")

	if err != nil {
		return nil, err
	}

	username, err := t.GetString("username")

	if err != nil {
		return nil, err
	}

	issuedAt, err := t.GetIssuedAt()

	if err != nil {
		return nil, err
	}

	expiredAt, err := t.GetExpiration()

	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        uuid.MustParse(id),
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}
