package api

import (
	"github.com/jaswdr/faker"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
	"github.com/kerok-kristoffer/formulating/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var f = faker.New()

func NewTestServer(t *testing.T, account db.UserAccount) *Server {
	config := util.Config{
		TokenSymmetricKey:   f.RandomStringWithLength(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, account)
	require.NoError(t, err)
	return server
}
