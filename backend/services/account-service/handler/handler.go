package handler

import (
	"encoding/json"
	"go-web-server/services/account-service/handler/middleware"
	"go-web-server/services/account-service/model"
	"go-web-server/services/account-service/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	service service.AccountService
}

func NewAccountHandler(service service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.HealthHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/accounts", h.CreateAccount)
		r.Get("/accounts/{accountId}", h.GetAccount)
		r.Post("/accounts/{accountId}/balance", h.UpdateBalance)
	})
}

func (h *AccountHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "UP"})
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID string `json:"customerId"`
		Currency   string `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Error decoding create account body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	acc, err := h.service.CreateAccount(r.Context(), body.CustomerID, body.Currency)
	if err != nil {
		log.Printf("Error creating account: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Account created: %s for customer %s", acc.ID, body.CustomerID)
	respondWithJSON(w, http.StatusCreated, acc)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountId")
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "Account ID is required")
		return
	}

	acc, err := h.service.GetAccount(r.Context(), accountID)
	if err != nil {
		log.Printf("Error getting account %s: %v", accountID, err)
		respondWithError(w, http.StatusNotFound, "Account not found")
		return
	}

	respondWithJSON(w, http.StatusOK, acc)
}

func (h *AccountHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountId")
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "Account ID is required")
		return
	}

	var body struct {
		Amount      float64               `json:"amount"`
		Type        model.LedgerEntryType `json:"type"`
		Description string                `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Error decoding update balance body for account %s: %v", accountID, err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.service.UpdateBalance(r.Context(), accountID, body.Amount, body.Type, body.Description)
	if err != nil {
		log.Printf("Error updating balance for account %s: %v", accountID, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Balance updated for account %s: %.2f (%s)", accountID, body.Amount, body.Type)
	w.WriteHeader(http.StatusNoContent)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}