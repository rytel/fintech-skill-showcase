package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go-web-server/database"
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
