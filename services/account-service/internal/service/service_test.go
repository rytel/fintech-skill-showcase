package service

import (
	"context"
	"testing"

	"go-web-server/services/account-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock of the AccountRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateAccount(acc *model.Account) error {
	args := m.Called(acc)
	return args.Error(0)
}

func (m *MockRepository) GetAccount(id string) (*model.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockRepository) UpdateBalance(accountID string, amount float64, entryType model.LedgerEntryType, description string) error {
	args := m.Called(accountID, amount, entryType, description)
	return args.Error(0)
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewAccountService(mockRepo)
	ctx := context.Background()

	customerID := uuid.New()
	currency := "USD"

	mockRepo.On("CreateAccount", mock.AnythingOfType("*model.Account")).Return(nil)

	acc, err := svc.CreateAccount(ctx, customerID.String(), currency)

	assert.NoError(t, err)
	assert.NotNil(t, acc)
	assert.Equal(t, currency, acc.Currency)
	assert.Equal(t, customerID, acc.CustomerID)
	mockRepo.AssertExpectations(t)
}

func TestGetAccount(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewAccountService(mockRepo)
	ctx := context.Background()

	accountID := uuid.New()
	expectedAcc := &model.Account{ID: accountID, Currency: "PLN"}

	mockRepo.On("GetAccount", accountID.String()).Return(expectedAcc, nil)

	acc, err := svc.GetAccount(ctx, accountID.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedAcc, acc)
	mockRepo.AssertExpectations(t)
}
