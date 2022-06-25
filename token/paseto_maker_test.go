package token

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {

	maker, err := NewJWTMaker(f.RandomStringWithLength(32))
	require.NoError(t, err)

	username := f.Internet().User()
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
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(f.RandomStringWithLength(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(f.Internet().User(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(f.RandomStringWithLength(32))
	require.NoError(t, err)
	maker2, err := NewPasetoMaker(f.RandomStringWithLength(32))
	require.NoError(t, err)

	duration := time.Minute
	token, err := maker.CreateToken(f.Internet().User(), duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	invalidToken, err := maker2.CreateToken(f.Internet().User(), duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(invalidToken) // todo kerok - add more examples of invalid tokens for test
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
