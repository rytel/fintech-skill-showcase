package main

import (
	"log"
	"net/http"
	"os"

	"go-web-server/api"
	"go-web-server/database"
)

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
	http.HandleFunc("/", api.SimpleTextHandler)     // Stary endpoint tekstowy
	http.HandleFunc("/api/status", api.JsonHandler) // NOWY: Endpoint JSON
	// Przekazujemy 'db' do handlera za pomocą domknięcia (closure)
	http.HandleFunc("/api/account/", api.MakeAccountHandler(db))
	// Rejestracja nowego endpointu transakcyjnego
	http.HandleFunc("/api/transactions", api.MakeTransactionHandler(db))

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