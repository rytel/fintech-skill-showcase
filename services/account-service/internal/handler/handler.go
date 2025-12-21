package handler

import (
	"encoding/json"
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
	// Stub
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	// Stub
}

func (h *AccountHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	// Stub
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
