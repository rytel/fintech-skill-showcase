package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-web-server/services/account-service/model"
	"go-web-server/services/account-service/service"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var jwtKey = []byte("my_secret_key_for_testing_only")

func createToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

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

func setupRouter(mockSvc service.AccountService) chi.Router {
	r := chi.NewRouter()
	h := NewAccountHandler(mockSvc)
	h.RegisterRoutes(r)
	return r
}

func TestCreateAccountHandler(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

	customerID := uuid.New().String()
	currency := "USD"
	reqBody, _ := json.Marshal(map[string]string{
		"customerId": customerID,
		"currency":   currency,
	})

	expectedAcc := &model.Account{ID: uuid.New(), Currency: currency}
	mockSvc.On("CreateAccount", mock.Anything, customerID, currency).Return(expectedAcc, nil)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+createToken("test_user"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var returnedAcc model.Account
	json.NewDecoder(rr.Body).Decode(&returnedAcc)
	assert.Equal(t, expectedAcc.ID, returnedAcc.ID)
}

func TestGetAccountHandler(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

	accountID := uuid.New()
	expectedAcc := &model.Account{ID: accountID, Currency: "PLN"}
	mockSvc.On("GetAccount", mock.Anything, accountID.String()).Return(expectedAcc, nil)

	req, _ := http.NewRequest("GET", "/accounts/"+accountID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+createToken("test_user"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var returnedAcc model.Account
	json.NewDecoder(rr.Body).Decode(&returnedAcc)
	assert.Equal(t, accountID, returnedAcc.ID)
}

func TestUpdateBalanceHandler(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

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
	req.Header.Set("Authorization", "Bearer "+createToken("test_user"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestAuthMiddleware_Failure(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

	req, _ := http.NewRequest("GET", "/accounts/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHealthHandler(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "UP")
}

func TestCreateAccountHandler_InvalidJSON(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBufferString("invalid json"))
	req.Header.Set("Authorization", "Bearer "+createToken("test_user"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAccountHandler_NotFound(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

	accountID := uuid.New().String()
	mockSvc.On("GetAccount", mock.Anything, accountID).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/accounts/"+accountID, nil)
	req.Header.Set("Authorization", "Bearer "+createToken("test_user"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateAccountHandler_ServiceError(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

	mockSvc.On("CreateAccount", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)

	reqBody, _ := json.Marshal(map[string]string{"customerId": "id", "currency": "USD"})
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+createToken("test_user"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateBalanceHandler_InvalidJSON(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

	req, _ := http.NewRequest("POST", "/accounts/some-id/balance", bytes.NewBufferString("invalid json"))
	req.Header.Set("Authorization", "Bearer "+createToken("test_user"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateBalanceHandler_ServiceError(t *testing.T) {
	mockSvc := new(MockService)
	r := setupRouter(mockSvc)

	mockSvc.On("UpdateBalance", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

	reqBody, _ := json.Marshal(map[string]interface{}{"amount": 10.0, "type": "DEPOSIT", "description": "desc"})
	req, _ := http.NewRequest("POST", "/accounts/some-id/balance", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+createToken("test_user"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}