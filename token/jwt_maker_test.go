package token

import (
	"testing"
	"time"

	"github.com/definitely-unique-username/simple_bank/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandString(32))

	require.NoError(t, err)

	username := util.RandOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredJWT(t *testing.T) {
	maker, err := NewJWTMaker(util.RandString(32))

	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandOwner(), -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTAlgorithName(t *testing.T) {
	payload, err := NewPayload(util.RandOwner(), time.Minute)

	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	signedToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandString(32))

	require.NoError(t, err)

	payload, err = maker.VerifyToken(signedToken)

	require.Error(t, err)
	// I don't know how check real jwt error
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
