package db

import (
	"context"
	"fmt"
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
	existed := make(map[int]bool)
	fmt.Println(">>before transactions: fromAccount balance: ", fromAccount.Balance, " toAccount balance: ", toAccount.Balance)

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
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, fromAccount.ID, result.FromAccount.ID)

		// check to account
		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, toAccount.ID, result.ToAccount.ID)

		//check accounts balance
		fmt.Println(">>tx: fromAccount balance: ", result.FromAccount.Balance, " toAccount balance: ", result.ToAccount.Balance)
		diff1 := fromAccount.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - toAccount.Balance
		fmt.Printf("diff1: %v diff2: %v\n", diff1, diff2)
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //amount, 2* amount, 3*amount .......n*amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check the final balances
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	fmt.Println(">>after all transactions: fromAccount balance: ", updatedFromAccount.Balance, " toAccount balance: ", updatedToAccount.Balance)
	require.Equal(t, updatedFromAccount.Balance, fromAccount.Balance-int64(n)*amount)
	require.Equal(t, updatedToAccount.Balance, toAccount.Balance+int64(n)*amount)

}

func TestTransferTxDeadLock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(10)
	store := NewStore(testDB)
	errs := make(chan error)
	fmt.Println(">>before transactions: fromAccount balance: ", account1.Balance, " toAccount balance: ", account2.Balance)

	// 10 concurrent transactions
	n := 10
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			arg := TransferTxParam{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			}

			_, err := store.TransferTx(context.Background(), arg)
			errs <- err
		}()
	}

	// validate results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//check the final balances
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>after all transactions: fromAccount balance: ", updatedFromAccount.Balance, " toAccount balance: ", updatedToAccount.Balance)
	require.Equal(t, updatedFromAccount.Balance, account1.Balance)
	require.Equal(t, updatedToAccount.Balance, account2.Balance)

}
