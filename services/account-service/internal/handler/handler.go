package handler

import (
	"encoding/json"
	"go-web-server/services/account-service/internal/model"
	"go-web-server/services/account-service/internal/service"
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
	r.Post("/accounts", h.CreateAccount)
	r.Get("/accounts/{accountId}", h.GetAccount)
	r.Post("/accounts/{accountId}/balance", h.UpdateBalance)
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID string `json:"customerId"`
		Currency   string `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	acc, err := h.service.CreateAccount(r.Context(), body.CustomerID, body.Currency)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

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
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.service.UpdateBalance(r.Context(), accountID, body.Amount, body.Type, body.Description)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

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