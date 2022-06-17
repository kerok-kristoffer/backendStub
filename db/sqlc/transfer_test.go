package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

// These are not really applicable to my current purposes,
// only followed tut for some example golang coding
func TestTransferTx(t *testing.T) {
	userAccount := NewUserAccount(testDB)

	user1 := CreateRandomUser(t)
	user2 := CreateRandomUser(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := userAccount.TransferTx(context.Background(), TransferTxParams{
				FromUserID: user1.ID,
				ToUserID:   user2.ID,
				Amount:     amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results
		require.NoError(t, err)
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, user1.ID, transfer.FromUserID.Int64)
		require.Equal(t, user2.ID, transfer.ToUserID.Int64)
		require.Equal(t, amount, transfer.Amount.Int64)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = userAccount.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		// Check entries
		require.NotEmpty(t, fromEntry)
		require.Equal(t, user1.ID, fromEntry.UserID.Int64)
		require.Equal(t, amount, fromEntry.Amount.Int64)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = userAccount.GetEntryByUserId(context.Background(), fromEntry.UserID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, user2.ID, toEntry.UserID.Int64)
		require.Equal(t, amount, toEntry.Amount.Int64)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = userAccount.GetEntryByUserId(context.Background(), toEntry.UserID)
		require.NoError(t, err)

		// TODO check account balances
		err = userAccount.DeleteTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		err = userAccount.DeleteEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		err = userAccount.DeleteEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
	}
	DeleteUser(t, user1)
	DeleteUser(t, user2)
}
