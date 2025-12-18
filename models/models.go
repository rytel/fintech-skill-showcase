package models

import "time"

// Account reprezentuje konto użytkownika w systemie.
type Account struct {
	ID        int     `json:"id"`
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// TransactionType definiuje typ transakcji (Deposit/Withdraw).
type TransactionType string

const (
	Deposit  TransactionType = "DEPOSIT"
	Withdraw TransactionType = "WITHDRAW"
)

// Transaction reprezentuje pojedynczą operację finansową na koncie.
type Transaction struct {
	ID        int             `json:"id"`
	AccountID int             `json:"account_id"`
	Type      TransactionType `json:"type"`
	Amount    float64         `json:"amount"`
	CreatedAt time.Time       `json:"created_at"`
}
