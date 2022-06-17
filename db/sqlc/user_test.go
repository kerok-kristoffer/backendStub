package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		FullName: F.Person().Name(),
		Hash:     F.Hash().MD5(),
		UserName: F.Internet().User(),
		Email:    F.Internet().Email(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Hash, user.Hash)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	err = testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}

func CreateRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		FullName: F.Person().Name(),
		Hash:     F.Hash().MD5(),
		UserName: F.Internet().User(),
		Email:    F.Internet().Email(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)

	return user
}

func TestGetUser(t *testing.T) {
	createdUser := CreateRandomUser(t)
	getUser, err := testQueries.GetUser(context.Background(), createdUser.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getUser)

	require.Equal(t, createdUser.ID, getUser.ID)
	require.Equal(t, createdUser.FullName, getUser.FullName)
	require.Equal(t, createdUser.UserName, getUser.UserName)
	require.Equal(t, createdUser.Email, getUser.Email)
	require.Equal(t, createdUser.Hash, getUser.Hash)
	require.WithinDurationf(t, createdUser.CreatedAt, getUser.CreatedAt, time.Second, "Error, created_at timestamps not within 1sec")

	DeleteUser(t, createdUser)
}

func DeleteUser(t *testing.T, user User) {
	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	createdUser := CreateRandomUser(t)
	getUser, err := testQueries.GetUserByUserEmail(context.Background(), createdUser.Email)

	require.NoError(t, err)
	require.NotEmpty(t, getUser)

	require.Equal(t, createdUser.ID, getUser.ID)
	require.Equal(t, createdUser.FullName, getUser.FullName)
	require.Equal(t, createdUser.UserName, getUser.UserName)
	require.Equal(t, createdUser.Email, getUser.Email)
	require.Equal(t, createdUser.Hash, getUser.Hash)
	require.WithinDurationf(t, createdUser.CreatedAt, getUser.CreatedAt, time.Second, "Error, created_at timestamps not within 1sec")

	err = testQueries.DeleteUser(context.Background(), createdUser.ID)
	require.NoError(t, err)
}

func TestGetUserByUserName(t *testing.T) {
	createdUser := CreateRandomUser(t)
	getUser, err := testQueries.GetUserByUserName(context.Background(), createdUser.UserName)

	require.NoError(t, err)
	require.NotEmpty(t, getUser)

	require.Equal(t, createdUser.ID, getUser.ID)
	require.Equal(t, createdUser.FullName, getUser.FullName)
	require.Equal(t, createdUser.UserName, getUser.UserName)
	require.Equal(t, createdUser.Email, getUser.Email)
	require.Equal(t, createdUser.Hash, getUser.Hash)
	require.WithinDurationf(t, createdUser.CreatedAt, getUser.CreatedAt, time.Second, "Error, created_at timestamps not within 1sec")

	err = testQueries.DeleteUser(context.Background(), createdUser.ID)
	require.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	createdUser := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), createdUser.ID)
	require.NoError(t, err)

	getUserExpectEmpty, err := testQueries.GetUser(context.Background(), createdUser.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getUserExpectEmpty)
}

func TestUpdateUserHash(t *testing.T) {
	// todo instead of creating new random user for each test, maybe use a constant
	// that is cleaned up after for all these "update" tests

	createdUser := CreateRandomUser(t)

	userHashParams := UpdateUserHashParams{
		ID:   createdUser.ID,
		Hash: F.Hash().MD5(),
	}

	updatedUser, err := testQueries.UpdateUserHash(context.Background(), userHashParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, createdUser.ID, updatedUser.ID)
	require.Equal(t, createdUser.FullName, updatedUser.FullName)
	require.Equal(t, userHashParams.Hash, updatedUser.Hash)
	require.WithinDurationf(t, createdUser.CreatedAt, updatedUser.CreatedAt, time.Second, "Error, created_at timestamps not within 1sec")

	err = testQueries.DeleteUser(context.Background(), createdUser.ID)
	require.NoError(t, err)
}

func TestUpdateUserName(t *testing.T) {
	createdUser := CreateRandomUser(t)

	userNameParams := UpdateUserNameParams{
		ID:       createdUser.ID,
		FullName: F.Person().Name(),
	}

	updatedUser, err := testQueries.UpdateUserName(context.Background(), userNameParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, createdUser.ID, updatedUser.ID)
	require.Equal(t, createdUser.Hash, updatedUser.Hash)
	require.Equal(t, userNameParams.FullName, updatedUser.FullName)
	require.WithinDurationf(t, createdUser.CreatedAt, updatedUser.CreatedAt, time.Second, "Error, created_at timestamps not within 1sec")

	err = testQueries.DeleteUser(context.Background(), createdUser.ID)
	require.NoError(t, err)
}

func TestListUsers(t *testing.T) {
	createdUsers := []User{}

	for i := 0; i < 10; i++ {
		createdUsers = append(createdUsers, CreateRandomUser(t))
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

	for _, user := range createdUsers {
		DeleteUser(t, user)
	}
}
