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
)

type Handler struct {
	repo *repository.PostgresRepository
}

func NewHandler(repo *repository.PostgresRepository) *Handler {
	return &Handler{repo: repo}
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

	account, err := h.repo.GetAccount(userID)
	if err != nil {
		log.Printf("DB Error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if account == nil {
		http.Error(w, "Account Not Found", http.StatusNotFound)
		return
	}

	h.sendJSON(w, http.StatusOK, account)
}

func (h *Handler) getTransactions(w http.ResponseWriter, r *http.Request, userID string) {
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

	updatedAccount, err := h.repo.CreateTransaction(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendJSON(w, http.StatusOK, updatedAccount)
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
