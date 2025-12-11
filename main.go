package main

import (
	"encoding/json" // NOWY: Pakiet do kodowania i dekodowania JSON.
	"fmt"
	"log"
	"net/http"
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
	// Dwa oddzielne endpointy:
	http.HandleFunc("/", simpleTextHandler)     // Stary endpoint tekstowy
	http.HandleFunc("/api/status", jsonHandler) // NOWY: Endpoint JSON

	const port = ":8080"
	log.Printf("Serwer jest gotowy i nasłuchuje na porcie %s. \n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
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
