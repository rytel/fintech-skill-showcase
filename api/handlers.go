package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go-web-server/database"
	"go-web-server/models"
)

// StatusResponse definiuje strukturę danych, która będzie serializowana do JSON.
type StatusResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// MakeAccountHandler tworzy handler, który ma dostęp do obiektu bazy danych.
func MakeAccountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Oczekujemy ścieżki /api/account/{user_id}
		// Pobieramy ID z URL
		userID := strings.TrimPrefix(r.URL.Path, "/api/account/")
		if userID == "" || userID == r.URL.Path {
			http.Error(w, "Brak User ID w ścieżce", http.StatusBadRequest)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Dozwolona tylko metoda GET", http.StatusMethodNotAllowed)
			return
		}

		// Pobierz konto z bazy
		account, err := database.GetAccount(db, userID)
		if err != nil {
			log.Printf("Błąd bazy danych: %v", err)
			http.Error(w, "Błąd serwera", http.StatusInternalServerError)
			return
		}

		if account == nil {
			http.Error(w, "Konto nie znalezione", http.StatusNotFound)
			return
		}

		// Zwróć JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(account); err != nil {
			log.Printf("Błąd kodowania JSON: %v", err)
		}
	}
}

// MakeTransactionHandler obsługuje tworzenie nowych transakcji (wpłaty/wypłaty).
func MakeTransactionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Dozwolona tylko metoda POST", http.StatusMethodNotAllowed)
			return
		}

		var req models.TransactionRequest
		// Dekodowanie JSON z body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Nieprawidłowy format JSON", http.StatusBadRequest)
			return
		}

		// Prosta walidacja
		if req.Amount <= 0 {
			http.Error(w, "Kwota musi być większa od zera", http.StatusBadRequest)
			return
		}
		if req.UserID == "" {
			http.Error(w, "Brak User ID", http.StatusBadRequest)
			return
		}
		if req.Type != models.Deposit && req.Type != models.Withdraw {
			http.Error(w, "Nieprawidłowy typ transakcji (oczekiwano: DEPOSIT lub WITHDRAW)", http.StatusBadRequest)
			return
		}

		// Wykonanie transakcji w bazie
		updatedAccount, err := database.CreateTransaction(db, req)
		if err != nil {
			log.Printf("Błąd przetwarzania transakcji: %v", err)
			// Rozróżnienie błędów biznesowych od systemowych
			if err.Error() == "niewystarczające środki na koncie" || err.Error() == "użytkownik nie posiada konta" {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, "Błąd serwera podczas przetwarzania transakcji", http.StatusInternalServerError)
			}
			return
		}

		// Sukces - zwróć zaktualizowane konto
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Lub 201 Created
		if err := json.NewEncoder(w).Encode(updatedAccount); err != nil {
			log.Printf("Błąd kodowania odpowiedzi JSON: %v", err)
		}
	}
}

// SimpleTextHandler obsługuje żądania GET dla /
func SimpleTextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Witajcie w Dockerze z Go! Żądano ścieżki: %s\n", r.URL.Path)
	log.Printf("Odebrano żądanie dla ścieżki: %s\n", r.URL.Path)
}

// JsonHandler obsługuje żądania dla /api/status i zwraca JSON
func JsonHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Definiowanie danych do zwrócenia
	data := StatusResponse{
		Status:  "OK",
		Service: "Go Docker Server",
		Version: "1.0",
	}

	// 2. Ustawienie nagłówka Content-Type na JSON!
	w.Header().Set("Content-Type", "application/json")

	// 3. Serializacja danych i wysłanie
	// json.NewEncoder(w) tworzy enkoder, który pisze dane JSON bezpośrednio do odpowiedzi HTTP (w).
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Błąd kodowania JSON: %v", err)
		return
	}
	log.Println("Wysłano odpowiedź JSON dla /api/status")
}