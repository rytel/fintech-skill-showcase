package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-web-server/services/account-service/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateAccount(ctx context.Context, customerID string, currency string) (*model.Account, error) {
	args := m.Called(ctx, customerID, currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockService) GetAccount(ctx context.Context, accountID string) (*model.Account, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockService) UpdateBalance(ctx context.Context, accountID string, amount float64, entryType model.LedgerEntryType, description string) error {
	args := m.Called(ctx, accountID, amount, entryType, description)
	return args.Error(0)
}

func TestCreateAccountHandler(t *testing.T) {
	mockSvc := new(MockService)
	h := NewAccountHandler(mockSvc)

	customerID := uuid.New().String()
	currency := "USD"
	reqBody, _ := json.Marshal(map[string]string{
		"customerId": customerID,
		"currency":   currency,
	})

	expectedAcc := &model.Account{ID: uuid.New(), Currency: currency}
	mockSvc.On("CreateAccount", mock.Anything, customerID, currency).Return(expectedAcc, nil)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateAccount)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var returnedAcc model.Account
	json.NewDecoder(rr.Body).Decode(&returnedAcc)
	assert.Equal(t, expectedAcc.ID, returnedAcc.ID)
	mockSvc.AssertExpectations(t)
}

func TestGetAccountHandler(t *testing.T) {
	mockSvc := new(MockService)
	h := NewAccountHandler(mockSvc)
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}", h.GetAccount)

	accountID := uuid.New()
	expectedAcc := &model.Account{ID: accountID, Currency: "PLN"}
	mockSvc.On("GetAccount", mock.Anything, accountID.String()).Return(expectedAcc, nil)

	req, _ := http.NewRequest("GET", "/accounts/"+accountID.String(), nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var returnedAcc model.Account
	json.NewDecoder(rr.Body).Decode(&returnedAcc)
	assert.Equal(t, accountID, returnedAcc.ID)
	mockSvc.AssertExpectations(t)
}
