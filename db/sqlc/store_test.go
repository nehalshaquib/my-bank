package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	store := NewStore(testDB)
	arg := TransferTxParam{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        10,
	}
	result, err := store.TransferTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, arg.FromAccountID, result.Transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, result.Transfer.ToAccountID)
	require.Equal(t, arg.Amount, result.Transfer.Amount)
	require.Equal(t, arg.FromAccountID, result.FromEntry.AccountID)
	require.Equal(t, arg.ToAccountID, result.ToEntry.AccountID)
	require.Equal(t, arg.Amount, -result.FromEntry.Amount)
	require.Equal(t, arg.Amount, result.ToEntry.Amount)

	// require.Equal(t, acc1.ID, result.FromAccount.ID)
	// require.Equal(t, acc2.ID, result.ToAccount.ID)
	// require.Equal(t, acc1.Balance, result.FromAccount.Balance-arg.Amount)
	// require.Equal(t, acc2.Balance, result.ToAccount.Balance+arg.Amount)
}
