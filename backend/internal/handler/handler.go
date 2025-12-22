package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"go-web-server/internal/model"
	"go-web-server/internal/repository"
	accModel "go-web-server/services/account-service/model"
	accService "go-web-server/services/account-service/service"
)

type Handler struct {
	repo       *repository.PostgresRepository
	accService accService.AccountService
}

func NewHandler(repo *repository.PostgresRepository, accService accService.AccountService) *Handler {
	return &Handler{
		repo:       repo,
		accService: accService,
	}
}

func (h *Handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "OK",
		"service": "Fintech API",
		"version": "2.0",
	}
	h.sendJSON(w, http.StatusOK, data)
}

func (h *Handler) AccountHandler(w http.ResponseWriter, r *http.Request) {
	pathSuffix := strings.TrimPrefix(r.URL.Path, "/api/account/")
	
	if strings.HasSuffix(pathSuffix, "/transactions") {
		h.getTransactions(w, r, strings.TrimSuffix(pathSuffix, "/transactions"))
		return
	}

	h.getAccount(w, r, pathSuffix)
}

func (h *Handler) getAccount(w http.ResponseWriter, r *http.Request, userID string) {
	if userID == "" {
		http.Error(w, "Missing User ID", http.StatusBadRequest)
		return
	}

	// Proxy to new AccountService
	// Note: We need a way to find the account for this user.
	// Since we are refactoring, we'll assume for now that we get the first account.
	// In a real migration, we would have a lookup table or more sophisticated logic.
	
	// For now, we'll try to find an account for this customer via the new repo
	// This is a temporary bridge.
	acc, err := h.accService.GetAccount(r.Context(), userID) // Assuming userID here maps to accountID for simplicity in this bridge
	if err != nil {
		log.Printf("Service Error: %v", err)
		http.Error(w, "Account Not Found", http.StatusNotFound)
		return
	}

	h.sendJSON(w, http.StatusOK, h.mapToMonolithAccount(acc))
}

func (h *Handler) getTransactions(w http.ResponseWriter, r *http.Request, userID string) {
	// For transactions, we would also proxy, but for this track, we'll stick to basic account operations
	// and mark this as potentially legacy or requiring update.
	transactions, err := h.repo.GetTransactions(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.sendJSON(w, http.StatusOK, transactions)
}

func (h *Handler) TransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Proxy to new AccountService for balance updates
	var entryType accModel.LedgerEntryType
	finalAmount := req.Amount

	if req.Type == model.Deposit {
		entryType = accModel.Deposit
	} else {
		entryType = accModel.Withdrawal
		finalAmount = -req.Amount
	}

	// We assume req.UserID is the accountID for this bridge
	err := h.accService.UpdateBalance(r.Context(), req.UserID, finalAmount, entryType, "Legacy Transaction Proxy")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get updated account to return
	acc, _ := h.accService.GetAccount(r.Context(), req.UserID)
	h.sendJSON(w, http.StatusOK, h.mapToMonolithAccount(acc))
}

func (h *Handler) mapToMonolithAccount(acc *accModel.Account) *model.Account {
	if acc == nil {
		return nil
	}
	// Simple hash-based ID conversion for legacy compatibility
	legacyID := int(acc.ID[0])<<24 | int(acc.ID[1])<<16 | int(acc.ID[2])<<8 | int(acc.ID[3])
	
	return &model.Account{
		ID:        legacyID,
		UserID:    acc.CustomerID.String(),
		Balance:   acc.Balance,
		CreatedAt: acc.CreatedAt,
	}
}

func (h *Handler) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.ResetDB(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test environment reset"))
}

func (h *Handler) sendJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// Auth Logic
var jwtKey = []byte("my_secret_key_for_testing_only")

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if creds.Username != "test_user" || creds.Password != "password123" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString(jwtKey)
	h.sendJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}
