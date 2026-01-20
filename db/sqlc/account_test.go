package db

import (
	"context"
	"database/sql"
	"github/hungphat2018/simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	getAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getAccount)

	require.Equal(t, account.ID, getAccount.ID)
	require.Equal(t, account.Balance, getAccount.Balance)
	require.Equal(t, account.Currency, getAccount.Currency)
	require.WithinDuration(t, account.CreatedAt, getAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	arg := AddAccountBalanceParams{
		ID:     account.ID,
		Amount: util.RandomMoney(),
	}

	UpdatedAccount, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, UpdatedAccount)

	require.Equal(t, account.ID, UpdatedAccount.ID)
	require.Equal(t, account.Balance+arg.Amount, UpdatedAccount.Balance)
	require.Equal(t, account.Currency, UpdatedAccount.Currency)
	require.WithinDuration(t, account.CreatedAt, UpdatedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, account2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
