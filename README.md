# Fintech Skill Showcase - Go Backend

Profesjonalny serwer backendowy napisany w Go, demonstrujący architekturę systemów finansowych, obsługę transakcji ACID oraz nowoczesny layout projektu.

## 🚀 Architektura i Struktura Projektu

Projekt stosuje **Standard Go Project Layout**, zapewniający czystą separację obaw:

- `cmd/server/`: Punkt wejścia aplikacji (minimalny bootstrap).
- `internal/app/`: Logika inicjalizacji i wstrzykiwania zależności (DI).
- `internal/handler/`: Warstwa transportowa (HTTP/REST), obsługa żądań i odpowiedzi (niezależna od bazy danych).
- `internal/repository/`: Warstwa dostępu do danych (PostgreSQL), implementacja logiki transakcyjnej (ACID, `FOR UPDATE`).
- `internal/model/`: Definicje struktur danych i modeli domenowych.
- `migrations/`: Skrypty SQL definiujące schemat bazy danych.
- `tests/`: Testy integracyjne (End-to-End).

## 🛠 Stos Technologiczny

- **Język:** Go 1.23
- **Baza danych:** PostgreSQL 15
- **Konteneryzacja:** Docker & Docker Compose (Multi-stage build)
- **Testowanie:** Go Testing Package + `sqlmock` dla testów jednostkowych

## 📜 Standardy i Kontrakty

### Standardy Kodowania
1.  **Dependency Injection**: Nie używamy zmiennych globalnych. Zależności są przekazywane przez konstruktory (np. `NewHandler`, `NewPostgresRepository`).
2.  **Enkapsulacja**: Cała kluczowa logika biznesowa MUSI znajdować się w `internal/`, aby nie była dostępna dla zewnętrznych modułów.
3.  **Bezpieczeństwo (Safety)**: Operacje na saldzie muszą być atomowe i wykonywane w ramach transakcji bazy danych.
4.  **Testowanie**:
    -   Testy jednostkowe z `sqlmock` dla repozytoriów.
    -   Testy integracyjne w `tests/` wymagające uruchomionego środowiska (Docker Compose).

### Kontrakt API
-   Wszystkie endpointy zwracają **JSON** z nagłówkiem `Content-Type: application/json`.
-   Błędy API są zwracane w czytelnym formacie z odpowiednimi kodami HTTP (400, 404, 500).

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
- **Automatyczne Migracje:** Serwer automatycznie inicjalizuje schemat bazy danych przy starcie.
- **Integracja iOS:** Dedykowane endpointy i wsparcie dla testów UI (zobacz [IOS_API.md](IOS_API.md)).