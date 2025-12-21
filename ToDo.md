# 🧭 Plan Rozwoju Aplikacji Finansowej (Go & Docker)

Lista zadań jest podzielona na fazy, co pozwala na stopniową implementację logiki biznesowej, infrastruktury i bezpieczeństwa.

---

## Faza 1: Fundamenty Serwera Go i Bazy Danych (Zakończona)

### 💾 Integracja z PostgreSQL (Lokalnie)

- [x] **1.1** Instalacja sterownika PostgreSQL w Go (`github.com/lib/pq`).
- [x] **1.2** Konfiguracja połączenia z bazą danych w kodzie Go (użycie zmiennych środowiskowych dla parametrów połączenia).
- [x] **1.3** Utworzenie struktur Go (Modeli): `Account` i `Transaction`.
- [x] **1.4** Zaimplementowanie narzędzia do migracji DB i utworzenie tabel.

### 🧪 Testy i Pierwsze API

- [x] **1.5** Utworzenie pierwszego testu jednostkowego (`*_test.go`).
- [x] **1.6** Implementacja endpointu **GET** `/api/account/{user_id}`.

---

## Faza 2: Transakcje i Logika Finansowa (Zakończona)

### 💰 Logika Transakcyjna (Kluczowe)

- [x] **2.1** Utworzenie endpointu **POST** `/api/transactions` (wpłaty/wypłaty).
- [x] **2.2** Obsługa JSON wejściowego (`TransactionRequest`).
- [x] **2.3** Zaimplementowanie walidacji biznesowej.
- [x] **2.4** Zaimplementowanie **Transakcji DB (ACID)** (`BEGIN`, `COMMIT`, `ROLLBACK`).

### 🎯 Testy Integracyjne

- [x] **2.5** Napisanie testów integracyjnych weryfikujących spójność danych.

---

## Faza 3: Docker Compose i Infrastruktura (Zakończona)

### 🐳 Lokalna Orkiestracja

- [x] **3.1** Utworzenie `docker-compose.yml` (`api-service` + `db-service`).
- [x] **3.2** Konfiguracja sieci między kontenerami.
- [x] **3.3** Wolumeny dla trwałości danych PostgreSQL.

---

## Faza 4: Modularny Monolit i Refaktoryzacja (Zakończona - Obecny Stan)

### 🏗 Nowoczesna Architektura (Modular Monolith)

- [x] **4.1** Wprowadzenie **Standard Go Project Layout** (`cmd/`, `internal/`).
- [x] **4.2** Wyodrębnienie domeny `AccountService` do `services/account-service/`.
- [x] **4.3** Migracja identyfikatorów z `int` na **UUID** (`google/uuid`).
- [x] **4.4** Implementacja wzorca **Gateway**: Główny handler (`internal/handler`) deleguje logikę do `AccountService`.
- [x] **4.5** Aktualizacja schematu bazy danych (wsparcie dla UUID i relacji customer-account).

### 🛡 Jakość Kodu

- [x] **4.6** Dependency Injection (brak zmiennych globalnych).
- [x] **4.7** Testy jednostkowe z `sqlmock`.
- [x] **4.8** Testy integracyjne zaktualizowane do nowej architektury i UUID.

---

## Faza 5: Aplikacja iOS i API Mobilne (W Trakcie)

### 🍎 Integracja Frontend

- [ ] **5.1** Utworzenie interfejsu SwiftUI (zewnętrzne repozytorium).
- [ ] **5.2** Klient sieciowy Swift (URLSession).
- [x] **5.3** Autoryzacja JWT (Backend Ready - `/api/login`).
- [ ] **5.4** Bezpieczne przechowywanie tokenów (Keychain).

### 🧪 Wsparcie Testów UI

- [x] **5.5** Endpoint `/api/test/reset` do resetowania środowiska.
- [x] **5.6** Endpoint `/api/account/{user_id}/transactions`.

---

## Faza 6: Pełna Architektura Mikroserwisów (Planowane)

### 🚀 Skalowanie i Rozproszenie

- [ ] **6.1** Wprowadzenie komunikacji gRPC/HTTP między Gateway a Account Service (zastąpienie importu pakietu).
- [ ] **6.2** Osobna baza danych dla Account Service.
- [ ] **6.3** Wydzielenie `account-service` do osobnego kontenera w Docker Compose.

---

## Faza 7: Wdrożenie w Chmurze (Planowane)

### ☁️ Produkcja

- [ ] **7.1** Cloud SQL (PostgreSQL).
- [ ] **7.2** Wypchnięcie obrazów do Container Registry.
- [ ] **7.3** Wdrożenie na platformę Serverless (Cloud Run).