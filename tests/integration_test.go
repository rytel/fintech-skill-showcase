package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"go-web-server/models"
)

const (
	apiURL = "http://localhost:8080"
	dbDSN  = "host=localhost port=5432 user=postgres password=secret dbname=fintech_db sslmode=disable"
)

// Helper to reset DB state for a specific user
func resetUser(t *testing.T, userID string) {
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Clean up transactions and account for the user
	// Note: In a real scenario, we might want to truncate or use a separate test DB.
	// For now, we cascade delete or delete by ID if we had it, but deleting by user_id requires knowing the account_id for transactions.
	// Let's simplified: delete account (which should cascade transactions if set up, but let's do it manually if not)
	
	// First get account ID
	var accountID int
	err = db.QueryRow("SELECT id FROM accounts WHERE user_id = $1", userID).Scan(&accountID)
	if err == nil {
		_, _ = db.Exec("DELETE FROM transactions WHERE account_id = $1", accountID)
		_, _ = db.Exec("DELETE FROM accounts WHERE id = $1", accountID)
	}
}

func TestTransactionFlow(t *testing.T) {
	// 0. Setup
	testUserID := "integration_test_user"
	resetUser(t, testUserID)
	defer resetUser(t, testUserID)

	// 1. Create Account (Simulated by inserting directly to DB, as we don't have CreateAccount API yet)
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO accounts (user_id, balance) VALUES ($1, $2)", testUserID, 0.0)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// 2. Test Deposit
	depositReq := models.TransactionRequest{
		UserID: testUserID,
		Type:   models.Deposit,
		Amount: 100.0,
	}
	payload, _ := json.Marshal(depositReq)

	resp, err := client.Post(apiURL+"/api/transactions", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to send deposit request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for deposit, got %d", resp.StatusCode)
	}

	var account models.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if account.Balance != 100.0 {
		t.Errorf("Expected balance 100.0 after deposit, got %f", account.Balance)
	}

	// 3. Test Withdraw
	withdrawReq := models.TransactionRequest{
		UserID: testUserID,
		Type:   models.Withdraw,
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
	overdraftReq := models.TransactionRequest{
		UserID: testUserID,
		Type:   models.Withdraw,
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
