package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	AccountActive AccountStatus = "active"
	AccountFrozen AccountStatus = "frozen"
	AccountClosed AccountStatus = "closed"
)

type Customer struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"externalId"`
	FullName   string    `json:"fullName"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Account struct {
	ID            uuid.UUID     `json:"id"`
	CustomerID    uuid.UUID     `json:"customerId"`
	AccountNumber string        `json:"accountNumber"`
	Currency      string        `json:"currency"`
	Balance       float64       `json:"balance"` // Using float64 for simplicity in example, but decimal is better for prod
	Status        AccountStatus `json:"status"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

type LedgerEntryType string

const (
	Deposit     LedgerEntryType = "deposit"
	Withdrawal  LedgerEntryType = "withdrawal"
	TransferIn  LedgerEntryType = "transfer_in"
	TransferOut LedgerEntryType = "transfer_out"
	FXExchange  LedgerEntryType = "fx_exchange"
)

type LedgerEntry struct {
	ID           uuid.UUID       `json:"id"`
	AccountID    uuid.UUID       `json:"accountId"`
	Type         LedgerEntryType `json:"type"`
	Amount       float64         `json:"amount"`
	BalanceAfter float64         `json:"balanceAfter"`
	ReferenceID  *uuid.UUID      `json:"referenceId,omitempty"`
	Description  string          `json:"description"`
	CreatedAt    time.Time       `json:"createdAt"`
}
