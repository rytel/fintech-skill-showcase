package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go-web-server/internal/model"
)

const (
	apiURL = "http://localhost:8080"
	dbDSN  = "host=localhost port=5432 user=postgres password=secret dbname=fintech_db sslmode=disable"
)

// Helper to reset DB state (Truncate all relevant tables)
func resetDB(t *testing.T) {
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Truncate tables to ensure clean state
	// Order matters due to foreign keys: ledger_entries -> accounts -> customers
	_, err = db.Exec("TRUNCATE TABLE ledger_entries, accounts, customers CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to truncate tables: %v", err)
	}
}

func TestTransactionFlow(t *testing.T) {
	// 0. Setup
	resetDB(t)
	// We don't defer reset here to allow inspection if test fails, 
	// but normally you might want to. Ideally use a unique ID per test run.

	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// 1. Create Customer and Account directly in DB
	customerID := uuid.New()
	accountID := uuid.New()
	accountNumber := "PL" + uuid.New().String()[0:20] // Fake IBAN-ish
	
	// Create Customer
	_, err = db.Exec("INSERT INTO customers (id, external_id, full_name) VALUES ($1, $2, $3)", 
		customerID, "ext_user_"+customerID.String(), "Integration Test User")
	if err != nil {
		t.Fatalf("Failed to create customer: %v", err)
	}

	// Create Account
	// Note: New schema has balance default 0, currency required
	_, err = db.Exec(`INSERT INTO accounts (id, customer_id, account_number, currency, balance, status) 
		VALUES ($1, $2, $3, $4, $5, $6)`, 
		accountID, customerID, accountNumber, "USD", 0.0, "active")
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// 2. Test Deposit
	// IMPORTANT: The current Handler acts as a bridge and treats the 'UserID' field in JSON as the AccountID UUID.
	depositReq := model.TransactionRequest{
		UserID: accountID.String(), 
		Type:   model.Deposit,
		Amount: 100.0,
	}
	payload, _ := json.Marshal(depositReq)

	resp, err := client.Post(apiURL+"/api/transactions", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to send deposit request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read body for debug
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		t.Fatalf("Expected status 200 for deposit, got %d. Body: %s", resp.StatusCode, buf.String())
	}

	var account model.Account
	// Re-decode response
	// Note: The handler returns a mapped legacy model, so balance should be plain float
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		// Try to read body again if decode fails (though body is consumed, usually need to peek)
		t.Fatalf("Failed to decode response: %v", err)
	}

	if account.Balance != 100.0 {
		t.Errorf("Expected balance 100.0 after deposit, got %f", account.Balance)
	}

	// 3. Test Withdraw
	withdrawReq := model.TransactionRequest{
		UserID: accountID.String(),
		Type:   model.Withdraw,
		Amount: 40.0,
	}
	payload, _ = json.Marshal(withdrawReq)

	resp, err = client.Post(apiURL+"/api/transactions", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to send withdraw request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for withdraw, got %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if account.Balance != 60.0 {
		t.Errorf("Expected balance 60.0 after withdrawal, got %f", account.Balance)
	}

	// 4. Test Overdraft (Withdraw more than balance)
	overdraftReq := model.TransactionRequest{
		UserID: accountID.String(),
		Type:   model.Withdraw,
		Amount: 100.0,
	}
	payload, _ = json.Marshal(overdraftReq)

	resp, err = client.Post(apiURL+"/api/transactions", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to send overdraft request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for overdraft, got %d", resp.StatusCode)
	}
}