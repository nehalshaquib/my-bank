package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	amount := int64(10)
	store := NewStore(testDB)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// 5 concurrent transactions
	n := 5
	for i := 0; i < n; i++ {
		go func() {
			arg := TransferTxParam{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			}

			result, err := store.TransferTx(context.Background(), arg)
			errs <- err
			results <- result
		}()
	}

	// validate results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfers
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, fromAccount.ID, result.Transfer.FromAccountID)
		require.Equal(t, toAccount.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		//check from entry
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, fromAccount.ID, result.FromEntry.AccountID)
		require.Equal(t, amount, -result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)

		//check to entry
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, toAccount.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)

		//check from account
		// require.NotEmpty(t, result.ToAccount)
		// require.Equal(t, fromAccount.Balance, result.FromAccount.Balance-10*5)

	}

	// require.Equal(t, acc1.ID, result.FromAccount.ID)
	// require.Equal(t, acc2.ID, result.ToAccount.ID)
	// require.Equal(t, acc1.Balance, result.FromAccount.Balance-arg.Amount)
	// require.Equal(t, acc2.Balance, result.ToAccount.Balance+arg.Amount)
}
