package repository

import (
	"testing"

	"go-web-server/services/account-service/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresAccountRepository(db)
	acc := &model.Account{
		ID:            uuid.New(),
		CustomerID:    uuid.New(),
		AccountNumber: "PL1234567890",
		Currency:      "PLN",
		Balance:       0.0,
		Status:        model.AccountActive,
	}

	mock.ExpectExec("INSERT INTO accounts").
		WithArgs(acc.ID, acc.CustomerID, acc.AccountNumber, acc.Currency, acc.Balance, acc.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateAccount(acc)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresAccountRepository(db)
	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "customer_id", "account_number", "currency", "balance", "status", "created_at", "updated_at"}).
		AddRow(id, uuid.New(), "PL123", "PLN", 100.0, "active", "2025-12-21", "2025-12-21")

	mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id =").
		WithArgs(id).
		WillReturnRows(rows)

	acc, err := repo.GetAccount(id.String())
	assert.NoError(t, err)
	assert.NotNil(t, acc)
	assert.Equal(t, id, acc.ID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateBalance_ACID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresAccountRepository(db)
	accountID := uuid.New()
	amount := 50.0

	mock.ExpectBegin()
	// Lock for update
	mock.ExpectQuery("SELECT balance FROM accounts WHERE id = (.+) FOR UPDATE").
		WithArgs(accountID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(100.0))

	// Update account balance
	mock.ExpectExec("UPDATE accounts SET balance = (.+) WHERE id =").
		WithArgs(150.0, accountID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Insert ledger entry
	mock.ExpectExec("INSERT INTO ledger_entries").
		WithArgs(sqlmock.AnyArg(), accountID, model.Deposit, amount, 150.0, nil, "Test deposit").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repo.UpdateBalance(accountID.String(), amount, model.Deposit, "Test deposit")
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
