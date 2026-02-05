package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	// "kasir-api/pkg/validator"
	// "github.com/go-playground/validator/v10"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) SummaryByDate(startDate string, endDate string) (*models.TransactionSummary, error) {
	return s.repo.GetSummaryByDate(startDate, endDate)
}
