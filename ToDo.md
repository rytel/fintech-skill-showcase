# 🧭 Plan Rozwoju Aplikacji Finansowej (Go & Docker)

Lista zadań jest podzielona na fazy, co pozwala na stopniową implementację logiki biznesowej, infrastruktury i bezpieczeństwa.

---

## Faza 1: Fundamenty Serwera Go i Bazy Danych (Lokalne)

### 💾 Integracja z PostgreSQL (Lokalnie)

- [x] **1.1** Instalacja sterownika PostgreSQL w Go (`github.com/lib/pq`).
- [x] **1.2** Konfiguracja połączenia z bazą danych w kodzie Go (użycie zmiennych środowiskowych dla parametrów połączenia).
- [x] **1.3** Utworzenie struktur Go (Modeli): `Account` (ID, UserID, Balance) i `Transaction` (ID, AccountID, Type, Amount).
- [x] **1.4** Zaimplementowanie narzędzia do migracji DB (np. `goose` lub `migrate`) i utworzenie tabel dla `Account` i `Transaction`.

### 🧪 Testy i Pierwsze API

- [x] **1.5** Utworzenie pierwszego pliku testowego (`*_test.go`) i napisanie prostego **testu jednostkowego** dla funkcji (np. walidacji kwoty).
- [x] **1.6** Implementacja endpointu **GET** `/api/account/{user_id}`, który pobiera i zwraca saldo konta danego użytkownika (użycie ścieżki URL do pobrania ID).

---

## Faza 2: Transakcje i Logika Finansowa

### 💰 Logika Transakcyjna (Kluczowe)

- [x] **2.1** Utworzenie endpointu **POST** `/api/transactions` do obsługi wpłat (`Deposit`) i wypłat (`Withdraw`).
- [x] **2.2** Obsługa JSON wejściowego (payload) i dekodowanie go do struktury Go (`TransactionRequest`).
- [x] **2.3** Zaimplementowanie mechanizmu **walidacji** (np. sprawdzenie wystarczającego salda dla wypłaty).
- [x] **2.4** Zaimplementowanie **Transakcji DB (ACID)**: użycie `db.Begin()`, `tx.Commit()`, `tx.Rollback()` w celu zagwarantowania spójności danych przy operacjach na saldzie i historii.

### 🎯 Testy Integracyjne

- [x] **2.5** Napisanie **testów integracyjnych**, które uruchamiają serwer, wysyłają żądania POST i weryfikują, czy zmiany w PostgreSQL są poprawne (weryfikacja salda i historii).

---

## Faza 3: Docker Compose i Środowisko Wieloserwisowe

### 🐳 Lokalna Orkiestracja

- [x] **3.1** Utworzenie pliku **`docker-compose.yml`** definiującego serwis `api-service` (Twój Go Server) i `db-service` (PostgreSQL). 
- [x] **3.2** Skonfigurowanie połączenia z DB w `api-service` tak, aby używało nazwy serwisu PostgreSQL (`db-service`) jako hosta.
- [x] **3.3** Konfiguracja wolumenu (volume) dla PostgreSQL, aby dane bazy przetrwały restarty kontenera.
- [x] **3.4** Weryfikacja uruchomienia za pomocą `docker compose up -d` i sprawdzenie, że oba serwisy działają i komunikują się.

---

## Faza 4: Wdrożenie w Chmurze (Platforma Serwerless)

### ☁️ Produkcja i Dostęp

- [ ] **4.1** Utworzenie zdalnej, zarządzanej bazy danych **PostgreSQL** w wybranej chmurze (np. GCP Cloud SQL / Azure Database / AWS RDS).
- [ ] **4.2** Zbudowanie obrazu Go i wypchnięcie go do **Rejestru Kontenerów** (np. Docker Hub, GCR, ECR).
- [ ] **4.3** Wdrożenie obrazu Go na zarządzanej platformie serwerless (GCP Cloud Run / Azure Container Apps / AWS Fargate), podając mu Zmienne Środowiskowe połączenia ze zdalną bazą danych.
- [ ] **4.4** Weryfikacja działania endpointu `/api/account/{user_id}` pod publicznym adresem URL.

---

## Faza 5: Aplikacja iOS (SwiftUI i Fintech)

### 🍎 Frontend i Bezpieczeństwo

- [ ] **5.1** Utworzenie interfejsu SwiftUI do wyświetlania salda i wykonania transakcji.
- [ ] **5.2** Implementacja logiki sieciowej (Swift `URLSession` / `Decodable`) do pobierania salda z API Go.
- [ ] **5.3** Dodanie tokena autoryzacji (np. JWT) do nagłówka każdego żądania HTTP w Swift.
- [ ] **5.4** Zaimplementowanie **Keychain Services** w iOS do bezpiecznego przechowywania tokena autoryzacyjnego.
- [ ] **5.5** (Opcjonalnie) Badanie integracji z **SwiftData** w celu lokalnego buforowania historii transakcji.

---