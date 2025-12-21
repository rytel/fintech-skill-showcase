package repository

import (
	"database/sql"
	"fmt"
	"go-web-server/services/account-service/model"

	"github.com/google/uuid"
)

type AccountRepository interface {
	CreateAccount(acc *model.Account) error
	GetAccount(id string) (*model.Account, error)
	UpdateBalance(accountID string, amount float64, entryType model.LedgerEntryType, description string) error
}

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) CreateAccount(acc *model.Account) error {
	query := `INSERT INTO accounts (id, customer_id, account_number, currency, balance, status) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, acc.ID, acc.CustomerID, acc.AccountNumber, acc.Currency, acc.Balance, acc.Status)
	return err
}

func (r *PostgresAccountRepository) GetAccount(id string) (*model.Account, error) {
	query := `SELECT id, customer_id, account_number, currency, balance, status, created_at, updated_at 
	          FROM accounts WHERE id = $1`
	var acc model.Account
	err := r.db.QueryRow(query, id).Scan(
		&acc.ID, &acc.CustomerID, &acc.AccountNumber, &acc.Currency, &acc.Balance, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}

func (r *PostgresAccountRepository) UpdateBalance(accountID string, amount float64, entryType model.LedgerEntryType, description string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Lock account for update to ensure ACID
	var currentBalance float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, accountID).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("could not find or lock account: %w", err)
	}

	newBalance := currentBalance + amount

	// 2. Update balance
	_, err = tx.Exec(`UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("could not update balance: %w", err)
	}

	// 3. Create ledger entry
	entryID := uuid.New()
	_, err = tx.Exec(`INSERT INTO ledger_entries (id, account_id, type, amount, balance_after, reference_id, description) 
	                  VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		entryID, accountID, entryType, amount, newBalance, nil, description)
	if err != nil {
		return fmt.Errorf("could not create ledger entry: %w", err)
	}

	return tx.Commit()
}