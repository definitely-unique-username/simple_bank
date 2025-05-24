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

	userId := util.RandInt(1, 10)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userId, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.UserID, userId)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredJWT(t *testing.T) {
	maker, err := NewJWTMaker(util.RandString(32))

	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandInt(1, 10), -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTAlgorithName(t *testing.T) {
	payload, err := NewPayload(util.RandInt(1, 10), time.Minute)

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
