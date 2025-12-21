package repository

import (
	"database/sql"
	"go-web-server/services/account-service/internal/model"
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
	return nil
}

func (r *PostgresAccountRepository) GetAccount(id string) (*model.Account, error) {
	return nil, nil
}

func (r *PostgresAccountRepository) UpdateBalance(accountID string, amount float64, entryType model.LedgerEntryType, description string) error {
	return nil
}
