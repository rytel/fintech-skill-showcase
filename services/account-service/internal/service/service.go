package service

import (
	"context"
	"go-web-server/services/account-service/internal/model"
	"go-web-server/services/account-service/internal/repository"
)

type AccountService interface {
	CreateAccount(ctx context.Context, customerID string, currency string) (*model.Account, error)
	GetAccount(ctx context.Context, accountID string) (*model.Account, error)
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(ctx context.Context, customerID string, currency string) (*model.Account, error) {
	return nil, nil
}

func (s *accountService) GetAccount(ctx context.Context, accountID string) (*model.Account, error) {
	return nil, nil
}
