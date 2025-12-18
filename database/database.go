package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"go-web-server/models"

	_ "github.com/lib/pq" // Sterownik PostgreSQL
)

// InitDB inicjalizuje połączenie z bazą danych PostgreSQL.
// Pobiera konfigurację ze zmiennych środowiskowych.
func InitDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Domyślne wartości, jeśli zmienne nie są ustawione (przydatne do lokalnego devu bez dockera, jeśli potrzeba)
	// W środowisku produkcyjnym/dockerowym powinny być zawsze ustawione.
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if dbname == "" {
		dbname = "fintech_db"
	}

	// Konstrukcja Connection String (DSN)
	// sslmode=disable jest używane lokalnie; na produkcji zalecane require/verify-full
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Printf("Łączenie z bazą danych: host=%s port=%s dbname=%s user=%s", host, port, dbname, user)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("błąd otwierania połączenia db: %w", err)
	}

	// Sprawdzenie czy połączenie jest faktycznie aktywne (ping)
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("błąd pingowania bazy danych: %w", err)
	}

	log.Println("Pomyślnie połączono z bazą danych!")
	return db, nil
}

// Migrate uruchamia skrypty SQL w celu utworzenia tabel.
// W prostym rozwiązaniu czytamy plik schema.sql i wykonujemy go.
func Migrate(db *sql.DB) error {
	// Wczytaj zawartość pliku SQL
	schema, err := os.ReadFile("migrations/schema.sql")
	if err != nil {
		return fmt.Errorf("nie udało się wczytać pliku migracji: %w", err)
	}

	// Wykonaj SQL
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("błąd podczas wykonywania migracji: %w", err)
	}

	log.Println("Migracja bazy danych zakończona sukcesem (tabele utworzone).")
	return nil
}

// GetAccount pobiera informacje o koncie dla danego użytkownika.
func GetAccount(db *sql.DB, userID string) (*models.Account, error) {
	query := `SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = $1`

	row := db.QueryRow(query, userID)

	var account models.Account
	err := row.Scan(&account.ID, &account.UserID, &account.Balance, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Konto nie istnieje
		}
		return nil, fmt.Errorf("błąd pobierania konta: %w", err)
	}

	return &account, nil
}

// CreateTransaction wykonuje operację finansową w ramach transakcji DB (ACID).
func CreateTransaction(db *sql.DB, req models.TransactionRequest) (*models.Account, error) {
	// 1. Rozpoczęcie transakcji
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("błąd rozpoczęcia transakcji: %w", err)
	}

	// Defer rollback w przypadku błędu (jeśli commit się uda, rollback nic nie zrobi)
	defer tx.Rollback()

	// 2. Pobranie konta z blokadą wiersza (FOR UPDATE)
	// To zapobiega wyścigom (race conditions), gdy dwa żądania próbują zmienić saldo jednocześnie.
	var account models.Account
	queryAccount := `SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = $1 FOR UPDATE`
	err = tx.QueryRow(queryAccount, req.UserID).Scan(&account.ID, &account.UserID, &account.Balance, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("użytkownik nie posiada konta")
		}
		return nil, fmt.Errorf("błąd pobierania konta: %w", err)
	}

	// 3. Walidacja i obliczenie nowego salda
	newBalance := account.Balance
	if req.Type == models.Deposit {
		newBalance += req.Amount
	} else if req.Type == models.Withdraw {
		if account.Balance < req.Amount {
			return nil, errors.New("niewystarczające środki na koncie")
		}
		newBalance -= req.Amount
	} else {
		return nil, errors.New("nieprawidłowy typ transakcji")
	}

	// 4. Aktualizacja salda konta
	updateQuery := `UPDATE accounts SET balance = $1 WHERE id = $2`
	_, err = tx.Exec(updateQuery, newBalance, account.ID)
	if err != nil {
		return nil, fmt.Errorf("błąd aktualizacji salda: %w", err)
	}

	// 5. Zapisanie historii transakcji
	insertQuery := `INSERT INTO transactions (account_id, type, amount) VALUES ($1, $2, $3)`
	_, err = tx.Exec(insertQuery, account.ID, req.Type, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("błąd zapisu historii transakcji: %w", err)
	}

	// 6. Zatwierdzenie transakcji (Commit)
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("błąd zatwierdzania transakcji: %w", err)
	}

	// Zaktualizuj strukturę do zwrócenia
	account.Balance = newBalance
	return &account, nil
}