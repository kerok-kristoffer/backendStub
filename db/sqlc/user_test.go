package db

import (
	"context"
	"database/sql"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {

	CreateRandomUser(t)
	// add test for hash not
}

func CreateRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		FullName: F.Person().Name(),
		Hash:     F.Hash().MD5(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Hash, user.Hash)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Hash, user2.Hash)
	require.WithinDurationf(t, user1.CreatedAt, user2.CreatedAt, time.Second, "Error, created_at timestamps not within 1sec")

}

func TestUpdateUser(t *testing.T) {

	user1 := CreateRandomUser(t)

	// todo Add test for update UserName after running migration, currently UpdateUserName does not return anything
	f := faker.New()
	arg := UpdateUserHashParams{
		ID:   user1.ID,
		Hash: f.Hash().MD5(),
	}

	user2, err := testQueries.UpdateUserHash(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, arg.Hash, user2.Hash)
	require.WithinDurationf(t, user1.CreatedAt, user2.CreatedAt, time.Second, "Error, created_at timestamps not within 1sec")
}

func TestDeleteUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)

}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomUser(t)
	}
	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
	for _, user := range users {
		require.NotEmpty(t, user)
	}

}