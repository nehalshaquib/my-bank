package db

import (
	"context"
	"testing"

	"github.com/nehalshaquib/my-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1, entry2)
}

func createEntryWithAccountId(t *testing.T, accountId int64) Entry {

	arg1 := CreateEntryParams{
		AccountID: accountId,
		Amount:    util.RandomMoney(),
	}
	entry1, err := testQueries.CreateEntry(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)
	require.Equal(t, accountId, entry1.AccountID)
	require.Equal(t, arg1.Amount, entry1.Amount)
	require.NotZero(t, entry1.ID)
	require.NotZero(t, entry1.CreatedAt)

	return entry1
}

func TestListEntries(t *testing.T) {
	acc := createRandomAccount(t)
	// create 10 random entries
	for i := 0; i < 10; i++ {
		createEntryWithAccountId(t, acc.ID)
	}

	arg := ListEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
