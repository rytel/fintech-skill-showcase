# Project Context: Fintech Skill Showcase

Ten dokument służy jako źródło prawdy o architekturze i standardach projektu dla deweloperów oraz agentów AI.

## 🏛 Architektura
Projekt realizuje **Standard Go Project Layout** z podziałem na:
- `cmd/server/`: Minimalny bootstrap.
- `internal/app/`: Inicjalizacja i wstrzykiwanie zależności (DI).
- `internal/handler/`: Warstwa HTTP (niezależna od bazy danych).
- `internal/repository/`: Warstwa SQL (logika transakcyjna ACID, `FOR UPDATE`).
- `internal/model/`: Modele domenowe.

## 🛠 Standardy Kodowania
1. **Dependency Injection**: Nie używamy zmiennych globalnych. Zależności są przekazywane przez konstruktory (`NewHandler`, `NewPostgresRepository`).
2. **Encapsulation**: Cała logika biznesowa MUSI znajdować się w `internal/`.
3. **Safety**: Operacje na saldzie muszą być atomowe (transakcje DB).
4. **Testing**:
   - Testy jednostkowe z `sqlmock` dla repozytoriów.
   - Testy integracyjne w `tests/` wymagające Docker Compose.

## 🚀 Infrastruktura
- **Docker**: Multi-stage build (`dockerfile`).
- **Compose**: `api-service` + `db-service` (PostgreSQL 15).
- **Go**: Wersja 1.23.

## 🔄 Kontrakt API
Zawsze zwracamy JSON z nagłówkiem `Content-Type: application/json`. Błędy powinny być czytelne i zwracać odpowiednie kody HTTP (400, 404, 500).
