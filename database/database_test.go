package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAccount(t *testing.T) {
	// 1. Tworzenie mocka bazy danych
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Stałe dane testowe
	userID := "user123"
	expectedBalance := 100.50
	expectedCreatedAt := time.Now()

	// SCENARIUSZ 1: Sukces - konto istnieje
	t.Run("powinien zwrócić konto gdy istnieje", func(t *testing.T) {
		// Oczekujemy zapytania SELECT.
		rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at"}).
			AddRow(1, userID, expectedBalance, expectedCreatedAt)

		mock.ExpectQuery(`^SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = \$1$`).
			WithArgs(userID).
			WillReturnRows(rows)

		// Wywołanie testowanej funkcji
		account, err := GetAccount(db, userID)
	
		// Weryfikacja wyników
		if err != nil {
			t.Errorf("nieoczekiwany błąd: %v", err)
		}
		if account == nil {
			t.Error("oczekiwano konta, otrzymano nil")
		}
		if account.UserID != userID {
			t.Errorf("nieprawidłowy UserID: oczekiwano %s, otrzymano %s", userID, account.UserID)
		}
		if account.Balance != expectedBalance {
			t.Errorf("nieprawidłowe saldo: oczekiwano %f, otrzymano %f", expectedBalance, account.Balance)
		}
	})

	// SCENARIUSZ 2: Błąd bazy danych
	t.Run("powinien zwrócić błąd gdy baza zwraca błąd", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = \$1$`).
			WithArgs(userID).
			WillReturnError(sql.ErrConnDone) // Symulujemy błąd połączenia

		account, err := GetAccount(db, userID)

		if err == nil {
			t.Error("oczekiwano błędu, ale go nie otrzymano")
		}
		if account != nil {
			t.Error("nie oczekiwano konta przy błędzie")
		}
	})

	// SCENARIUSZ 3: Konto nie istnieje (brak wyników)
	t.Run("powinien zwrócić nil gdy konto nie istnieje", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = \$1$`).
			WithArgs("unknown_user").
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at"})) // Pusty wynik

		account, err := GetAccount(db, "unknown_user")

		if err != nil {
			t.Errorf("nieoczekiwany błąd przy braku wyników: %v", err)
		}
		if account != nil {
			t.Errorf("oczekiwano nil, otrzymano: %+v", account)
		}
	})

	// Sprawdzenie, czy wszystkie oczekiwania wobec mocka zostały spełnione
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("nie spełniono oczekiwań mocka: %s", err)
	}
}
