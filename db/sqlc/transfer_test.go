package db

import (
	"context"
	"testing"

	"github.com/nehalshaquib/my-bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	transferResponse, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transferResponse)
	require.Equal(t, transferResponse, transfer)
}

func TestListTransfers(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	arg := ListTransfersParams{
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		Limit:         1,
		Offset:        0,
	}
	tranfersList, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tranfersList)
	require.Equal(t, tranfersList[0], transfer)
}

func TestListTransfersBetween(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	arg := ListTransfersBetweenParams{
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		Limit:         1,
		Offset:        0,
	}
	tranfersList, err := testQueries.ListTransfersBetween(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tranfersList)
	require.Equal(t, tranfersList[0], transfer)
}

func TestListTransfersFrom(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	arg := ListTransfersFromParams{
		FromAccountID: transfer.FromAccountID,
		Limit:         1,
		Offset:        0,
	}
	tranfersList, err := testQueries.ListTransfersFrom(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tranfersList)
	require.Equal(t, tranfersList[0], transfer)
}

func TestListTransfersTo(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	arg := ListTransfersToParams{
		ToAccountID: transfer.ToAccountID,
		Limit:       1,
		Offset:      0,
	}
	tranfersList, err := testQueries.ListTransfersTo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tranfersList)
	require.Equal(t, tranfersList[0], transfer)
}
