package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go-web-server/services/account-service/internal/handler"
	"go-web-server/services/account-service/internal/model"
	"go-web-server/services/account-service/internal/repository"
	"go-web-server/services/account-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDB *sql.DB
	jwtKey = []byte("my_secret_key_for_testing_only")
)

func createToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func TestMain(m *testing.M) {
	host := os.Getenv("DB_HOST")
	if host == "" { host = "localhost" }
	port := os.Getenv("DB_PORT")
	if port == "" { port = "5432" }
	user := os.Getenv("DB_USER")
	if user == "" { user = "postgres" }
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	if dbname == "" { dbname = "fintech_db_test" } // Use test DB

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	testDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Failed to connect to test DB: %v\n", err)
		os.Exit(1)
	}

	if err = testDB.Ping(); err != nil {
		fmt.Printf("Failed to ping test DB: %v\n", err)
		// We don't exit here to allow skipping tests if DB is not available
	}

	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func setupIntegration(t *testing.T) (chi.Router, *sql.DB) {
	if err := testDB.Ping(); err != nil {
		t.Skip("Database not available for integration tests")
	}

	// Clean up and Migrate
	_, err := testDB.Exec("DROP TABLE IF EXISTS ledger_entries, accounts, customers CASCADE")
	require.NoError(t, err)

	schema, err := os.ReadFile("../migrations/schema.sql")
	require.NoError(t, err)
	_, err = testDB.Exec(string(schema))
	require.NoError(t, err)

	repo := repository.NewPostgresAccountRepository(testDB)
	svc := service.NewAccountService(repo)
	h := handler.NewAccountHandler(svc)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	return r, testDB
}

func TestAccountLifecycle_Integration(t *testing.T) {
	r, db := setupIntegration(t)
	token := createToken("test_user")

	// 1. Create a Customer in DB (Manual step since we don't have Customer API yet)
	customerID := uuid.New()
	externalID := "auth_user_123"
	_, err := db.Exec("INSERT INTO customers (id, external_id, full_name) VALUES ($1, $2, $3)",
		customerID, externalID, "Integration Test User")
	require.NoError(t, err)

	// 2. Create Account via API
	createBody, _ := json.Marshal(map[string]string{
		"customerId": customerID.String(),
		"currency":   "USD",
	})
	req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(createBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var acc model.Account
	err = json.Unmarshal(rr.Body.Bytes(), &acc)
	require.NoError(t, err)
	assert.Equal(t, "USD", acc.Currency)
	assert.Equal(t, 0.0, acc.Balance)

	accountID := acc.ID.String()

	// 3. Get Account via API
	req = httptest.NewRequest("GET", "/accounts/"+accountID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &acc)
	require.NoError(t, err)
	assert.Equal(t, accountID, acc.ID.String())

	// 4. Update Balance (Deposit)
	depositBody, _ := json.Marshal(map[string]interface{}{
		"amount":      100.50,
		"type":        "DEPOSIT",
		"description": "Initial deposit",
	})
	req = httptest.NewRequest("POST", "/accounts/"+accountID+"/balance", bytes.NewBuffer(depositBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// 5. Verify Balance
	req = httptest.NewRequest("GET", "/accounts/"+accountID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &acc)
	assert.Equal(t, 100.50, acc.Balance)

	// 6. Update Balance (Withdraw - Success)
	withdrawBody, _ := json.Marshal(map[string]interface{}{
		"amount":      -50.00,
		"type":        "WITHDRAWAL",
		"description": "ATM withdrawal",
	})
	req = httptest.NewRequest("POST", "/accounts/"+accountID+"/balance", bytes.NewBuffer(withdrawBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// 7. Update Balance (Withdraw - Insufficient Funds)
	withdrawBodyFail, _ := json.Marshal(map[string]interface{}{
		"amount":      -1000.00,
		"type":        "WITHDRAWAL",
		"description": "Too much",
	})
	req = httptest.NewRequest("POST", "/accounts/"+accountID+"/balance", bytes.NewBuffer(withdrawBodyFail))
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "insufficient funds")
}
