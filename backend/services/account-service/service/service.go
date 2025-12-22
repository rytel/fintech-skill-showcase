package service

import (
	"context"
	"fmt"
	"go-web-server/services/account-service/model"
	"go-web-server/services/account-service/repository"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type AccountService interface {
	CreateAccount(ctx context.Context, customerID string, currency string) (*model.Account, error)
	GetAccount(ctx context.Context, accountID string) (*model.Account, error)
	UpdateBalance(ctx context.Context, accountID string, amount float64, entryType model.LedgerEntryType, description string) error
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(ctx context.Context, customerID string, currency string) (*model.Account, error) {
	custUUID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer id: %w", err)
	}

	acc := &model.Account{
		ID:            uuid.New(),
		CustomerID:    custUUID,
		AccountNumber: generateAccountNumber(),
		Currency:      currency,
		Balance:       0.0,
		Status:        model.AccountActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.CreateAccount(acc); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return acc, nil
}

func (s *accountService) GetAccount(ctx context.Context, accountID string) (*model.Account, error) {
	acc, err := s.repo.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if acc == nil {
		return nil, fmt.Errorf("account not found")
	}
	return acc, nil
}

func (s *accountService) UpdateBalance(ctx context.Context, accountID string, amount float64, entryType model.LedgerEntryType, description string) error {
	acc, err := s.repo.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if acc == nil {
		return fmt.Errorf("account not found")
	}

	// Business Rule: Ensure sufficient funds for withdrawals
	if amount < 0 && acc.Balance+amount < 0 {
		return fmt.Errorf("insufficient funds: current balance %.2f, requested withdrawal %.2f", acc.Balance, -amount)
	}

	if err := s.repo.UpdateBalance(accountID, amount, entryType, description); err != nil {
		return fmt.Errorf("failed to update balance in repository: %w", err)
	}

	return nil
}

// generateAccountNumber creates a dummy bank account number
// In a real system, this would follow IBAN or other standards
func generateAccountNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("PL%010d", r.Int63n(10000000000))
}