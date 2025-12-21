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

func TestUpdateBalanceHandler(t *testing.T) {
	mockSvc := new(MockService)
	h := NewAccountHandler(mockSvc)
	r := chi.NewRouter()
	r.Post("/accounts/{accountId}/balance", h.UpdateBalance)

	accountID := uuid.New().String()
	amount := 100.0
	entryType := model.Deposit
	description := "Test deposit"

	reqBody, _ := json.Marshal(map[string]interface{}{
		"amount":      amount,
		"type":        entryType,
		"description": description,
	})

	mockSvc.On("UpdateBalance", mock.Anything, accountID, amount, entryType, description).Return(nil)

	req, _ := http.NewRequest("POST", "/accounts/"+accountID+"/balance", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreateAccountHandler_InvalidJSON(t *testing.T) {
	mockSvc := new(MockService)
	h := NewAccountHandler(mockSvc)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateAccount)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAccountHandler_NotFound(t *testing.T) {
	mockSvc := new(MockService)
	h := NewAccountHandler(mockSvc)
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}", h.GetAccount)

	accountID := uuid.New().String()
	mockSvc.On("GetAccount", mock.Anything, accountID).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/accounts/"+accountID, nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateAccountHandler_ServiceError(t *testing.T) {
	mockSvc := new(MockService)
	h := NewAccountHandler(mockSvc)

	mockSvc.On("CreateAccount", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)

	reqBody, _ := json.Marshal(map[string]string{"customerId": "id", "currency": "USD"})
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateAccount)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateBalanceHandler_InvalidJSON(t *testing.T) {
	mockSvc := new(MockService)
	h := NewAccountHandler(mockSvc)
	r := chi.NewRouter()
	r.Post("/accounts/{accountId}/balance", h.UpdateBalance)

	req, _ := http.NewRequest("POST", "/accounts/some-id/balance", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateBalanceHandler_ServiceError(t *testing.T) {
	mockSvc := new(MockService)
	h := NewAccountHandler(mockSvc)
	r := chi.NewRouter()
	r.Post("/accounts/{accountId}/balance", h.UpdateBalance)

	mockSvc.On("UpdateBalance", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

	reqBody, _ := json.Marshal(map[string]interface{}{"amount": 10.0, "type": "DEPOSIT", "description": "desc"})
	req, _ := http.NewRequest("POST", "/accounts/some-id/balance", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
