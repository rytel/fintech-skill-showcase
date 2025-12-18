package main

import (
	"database/sql"
	"encoding/json" // NOWY: Pakiet do kodowania i dekodowania JSON.
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"go-web-server/database" // Importujemy nasz pakiet bazy danych
)

// Definiujemy strukturę danych, która będzie serializowana do JSON.
// Pola muszą być eksportowane (zaczynać się wielką literą), aby były widoczne dla pakietu json.
// Tag `json:"nazwa"` definiuje nazwę klucza w wynikowym JSON.
type StatusResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func main() {
	// 1. Inicjalizacja Bazy Danych
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Nie udało się połączyć z bazą danych: %v", err)
	}
	defer db.Close() // Zamknij połączenie, gdy funkcja main się zakończy (czyli przy zamknięciu serwera)

	// 2. Uruchomienie Migracji (utworzenie tabel)
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Nie udało się przeprowadzić migracji: %v", err)
	}

	// Dwa oddzielne endpointy:
	http.HandleFunc("/", simpleTextHandler)     // Stary endpoint tekstowy
	http.HandleFunc("/api/status", jsonHandler) // NOWY: Endpoint JSON
	// Przekazujemy 'db' do handlera za pomocą domknięcia (closure)
	http.HandleFunc("/api/account/", makeAccountHandler(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port
	log.Printf("Serwer jest gotowy i nasłuchuje na porcie %s. \n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

// makeAccountHandler tworzy handler, który ma dostęp do obiektu bazy danych.
func makeAccountHandler(db *sql.DB) http.HandlerFunc {
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

// Obsługuje żądania GET dla /
func simpleTextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Witajcie w Dockerze z Go! Żądano ścieżki: %s\n", r.URL.Path)
	log.Printf("Odebrano żądanie dla ścieżki: %s\n", r.URL.Path)
}

// NOWY: Obsługuje żądania dla /api/status i zwraca JSON
func jsonHandler(w http.ResponseWriter, r *http.Request) {
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
