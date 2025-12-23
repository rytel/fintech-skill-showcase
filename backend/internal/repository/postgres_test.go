package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"go-web-server/internal/model"
)

func TestGetAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at"}).
		AddRow("1", "test_user", 100.0, now)

	mock.ExpectQuery(`SELECT id, user_id, balance, created_at FROM accounts`).
		WithArgs("test_user").
		WillReturnRows(rows)

	account, err := repo.GetAccount("test_user")
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	if account.UserID != "test_user" || account.Balance != 100.0 {
		t.Errorf("expected test_user with 100.0 balance, got %v", account)
	}
}

func TestCreateTransaction_Deposit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	req := model.TransactionRequest{
		UserID: "test_user",
		Type:   model.Deposit,
		Amount: 50.0,
	}

	mock.ExpectBegin()
	// Lock account
	mock.ExpectQuery(`SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = \$1 FOR UPDATE`).
		WithArgs("test_user").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at"}).AddRow("1", "test_user", 100.0, time.Now()))
	
	// Update balance
	mock.ExpectExec(`UPDATE accounts SET balance = \$1 WHERE id = \$2`).
		WithArgs(150.0, "1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Insert transaction record - Note: The code uses 'transactions' table in CreateTransaction
	// but the GetTransactionsRaw uses 'ledger_entries'. 
	// I should check repository/postgres.go CreateTransaction implementation again.
	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs("1", model.Deposit, 50.0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	mock.ExpectCommit()

	account, err := repo.CreateTransaction(req)
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	if account.Balance != 150.0 {
		t.Errorf("expected balance 150.0, got %f", account.Balance)
	}
}

func TestCreateTransaction_Withdraw_InsufficientFunds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	req := model.TransactionRequest{
		UserID: "test_user",
		Type:   model.Withdraw,
		Amount: 200.0,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = \$1 FOR UPDATE`).
		WithArgs("test_user").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at"}).AddRow("1", "test_user", 100.0, time.Now()))
	
	mock.ExpectRollback()

	_, err = repo.CreateTransaction(req)
	if err == nil || err.Error() != "niewystarczające środki na koncie" {
		t.Errorf("expected insufficient funds error, got: %v", err)
	}
}

func TestGetTransactionsRaw(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	rows := sqlmock.NewRows([]string{"id", "account_id", "type", "amount", "created_at"}).
		AddRow("1", "1", "deposit", 100.0, time.Now()).
		AddRow("2", "1", "withdrawal", 50.0, time.Now())

	mock.ExpectQuery(`SELECT id, account_id, type, amount, created_at FROM ledger_entries WHERE account_id = \$1`).
		WithArgs("1").
		WillReturnRows(rows)

	txs, err := repo.GetTransactionsRaw("1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(txs) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(txs))
	}
}