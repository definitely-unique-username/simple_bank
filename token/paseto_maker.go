package token

import (
	"strconv"
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

func (m *PasetoMaker) CreateToken(userID int64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, duration)

	if err != nil {
		return "", payload, err
	}

	token := paseto.NewToken()

	token.SetString("id", payload.ID.String())
	token.SetString("userId", strconv.FormatInt(payload.UserID, 10))
	token.SetIssuedAt(payload.IssuedAt)
	token.SetNotBefore(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	return token.V4Encrypt(m.symmetricalKey, m.implicit), payload, nil
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

	userID, err := t.GetString("userId")

	if err != nil {
		return nil, err
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
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
		UserID:    userIDInt,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}
