package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandString(8)

	hash, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = VerifyPassword(password, hash)
	require.NoError(t, err)
}

func TextInvalidPassword(t *testing.T) {
	password := "passowrd"

	hash, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = VerifyPassword(password, "invalid_password")
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

func TestSamePasswords(t *testing.T) {
	password := "passowrd"

	hash1, err := HashPassword(password)
	require.NoError(t, err)

	hash2, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEqual(t, hash1, hash2)
}
