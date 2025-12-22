package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"go-web-server/internal/model"

	_ "github.com/lib/pq" // Sterownik PostgreSQL
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// InitDB inicjalizuje połączenie z bazą danych PostgreSQL.
func InitDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" { host = "localhost" }
	if port == "" { port = "5432" }
	if user == "" { user = "postgres" }
	if dbname == "" { dbname = "fintech_db" }

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("błąd otwierania połączenia db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("błąd pingowania bazy danych: %w", err)
	}

	return db, nil
}

func (r *PostgresRepository) Migrate() error {
	schema, err := os.ReadFile("migrations/schema.sql")
	if err != nil {
		return fmt.Errorf("nie udało się wczytać pliku migracji: %w", err)
	}

	_, err = r.db.Exec(string(schema))
	return err
}

func (r *PostgresRepository) GetAccount(userID string) (*model.Account, error) {
	query := `SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = $1`
	row := r.db.QueryRow(query, userID)

	var account model.Account
	err := row.Scan(&account.ID, &account.UserID, &account.Balance, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *PostgresRepository) CreateTransaction(req model.TransactionRequest) (*model.Account, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var account model.Account
	queryAccount := `SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = $1 FOR UPDATE`
	err = tx.QueryRow(queryAccount, req.UserID).Scan(&account.ID, &account.UserID, &account.Balance, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("użytkownik nie posiada konta")
		}
		return nil, err
	}

	newBalance := account.Balance
	if req.Type == model.Deposit {
		newBalance += req.Amount
	} else if req.Type == model.Withdraw {
		if account.Balance < req.Amount {
			return nil, errors.New("niewystarczające środki na koncie")
		}
		newBalance -= req.Amount
	}

	_, err = tx.Exec(`UPDATE accounts SET balance = $1 WHERE id = $2`, newBalance, account.ID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`INSERT INTO transactions (account_id, type, amount) VALUES ($1, $2, $3)`, account.ID, req.Type, req.Amount)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	account.Balance = newBalance
	return &account, nil
}

func (r *PostgresRepository) GetTransactions(userID string) ([]model.Transaction, error) {
	var accountID int
	err := r.db.QueryRow("SELECT id FROM accounts WHERE user_id = $1", userID).Scan(&accountID)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`SELECT id, account_id, type, amount, created_at FROM transactions WHERE account_id = $1 ORDER BY created_at DESC`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.AccountID, &t.Type, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *PostgresRepository) ResetDB() error {
	_, err := r.db.Exec("TRUNCATE TABLE transactions, accounts RESTART IDENTITY CASCADE")
	if err != nil {
		return err
	}
	_, err = r.db.Exec("INSERT INTO accounts (user_id, balance) VALUES ($1, $2)", "test_user", 1000.0)
	return err
}
