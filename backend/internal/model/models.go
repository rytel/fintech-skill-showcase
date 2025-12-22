package model

import "time"

// Account reprezentuje konto użytkownika w systemie.
type Account struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// TransactionType definiuje typ transakcji (Deposit/Withdraw).
type TransactionType string

const (
	Deposit  TransactionType = "DEPOSIT"
	Withdraw TransactionType = "WITHDRAWAL"
)

// Transaction reprezentuje pojedynczą operację finansową na koncie.
type Transaction struct {
	ID        string          `json:"id"`
	AccountID string          `json:"account_id"`
	Type      TransactionType `json:"type"`
	Amount    float64         `json:"amount"`
	CreatedAt time.Time       `json:"created_at"`
}

// TransactionRequest reprezentuje dane wejściowe dla nowej transakcji.
type TransactionRequest struct {
	UserID string          `json:"user_id"`
	Type   TransactionType `json:"type"`
	Amount float64         `json:"amount"`
}