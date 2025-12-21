package app

import (
	"log"
	"net/http"
	"os"

	"go-web-server/internal/handler"
	"go-web-server/internal/repository"
	accRepo "go-web-server/services/account-service/repository"
	accService "go-web-server/services/account-service/service"
)

func Run() {
	db, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresRepository(db)
	if err := repo.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Microservices integration
	newAccRepo := accRepo.NewPostgresAccountRepository(db)
	newAccService := accService.NewAccountService(newAccRepo)

	h := handler.NewHandler(repo, newAccService)
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/api/status", h.StatusHandler)
	mux.HandleFunc("/api/account/", h.AccountHandler)
	mux.HandleFunc("/api/transactions", h.TransactionHandler)
	mux.HandleFunc("/api/login", h.LoginHandler)
	mux.HandleFunc("/api/test/reset", h.ResetHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
