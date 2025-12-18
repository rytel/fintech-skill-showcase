# Fintech Skill Showcase - Go Backend

Profesjonalny serwer backendowy napisany w Go, demonstrujący architekturę systemów finansowych, obsługę transakcji ACID oraz nowoczesny layout projektu.

## 🚀 Architektura i Struktura Projektu

Projekt stosuje **Standard Go Project Layout**, zapewniający czystą separację obaw:

- `cmd/server/`: Punkt wejścia aplikacji.
- `internal/app/`: Logika inicjalizacji i bootstrapu serwera.
- `internal/handler/`: Warstwa transportowa (HTTP/REST), obsługa żądań i odpowiedzi.
- `internal/repository/`: Warstwa dostępu do danych (PostgreSQL), implementacja logiki transakcyjnej.
- `internal/model/`: Definicje struktur danych i modeli domenowych.
- `migrations/`: Skrypty SQL definiujące schemat bazy danych.
- `tests/`: Testy integracyjne (End-to-End).

## 🛠 Stos Technologiczny

- **Język:** Go 1.23
- **Baza danych:** PostgreSQL 15
- **Konteneryzacja:** Docker & Docker Compose
- **Testowanie:** Go Testing Package + `sqlmock` dla testów jednostkowych

## 🏁 Jak uruchomić?

### Szybki start (Docker)
Wymaga zainstalowanego Dockera i Docker Compose.

```bash
docker-compose up --build
```
Serwer będzie dostępny pod adresem: `http://localhost:8080`

### Uruchamianie Testów

**Testy jednostkowe (izolowane):**
```bash
go test ./internal/...
```

**Testy integracyjne (wymagają bazy danych):**
```bash
docker-compose up -d
go test -v ./tests/integration_test.go
```

## 🔐 Kluczowe Funkcjonalności

- **Transakcyjność ACID:** Wszystkie operacje finansowe (wpłaty/wypłaty) są wykonywane w ramach transakcji DB z blokowaniem wierszy (`FOR UPDATE`), co zapobiega wyścigom (race conditions).
- **Enkapsulacja:** Kluczowa logika biznesowa znajduje się w katalogu `internal/`, co gwarantuje, że nie zostanie ona użyta niezgodnie z przeznaczeniem przez zewnętrzne moduły.
- **Automatyczne Migracje:** Serwer automatycznie inicjalizuje schemat bazy danych przy starcie.
