package services

import (
	"bank-system/src/entities"
	"bank-system/src/repositories"
	"context"

	"github.com/sirupsen/logrus"
)

type TransactionService interface {
	Create(ctx context.Context, userId string, data entities.Transaction) (string, error)
	GetByAccountId(ctx context.Context, userId string, accountId string) ([]entities.Transaction, error)
}

type transactionService struct {
	transactionRepository repositories.TransactionRepository
	logger                *logrus.Logger
}

func NewTransactionService(
	transactionRepository repositories.TransactionRepository,
	logger *logrus.Logger,
) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		logger:                logger,
	}
}

func (this *transactionService) Create(ctx context.Context, userId string, data entities.Transaction) (string, error) {
	_, err := this.transactionRepository.Create(ctx, data)
	if err != nil {
		this.logger.Errorf("Failed to create transaction: %v", err)
		return "", err
	}

	return data.ID, nil
}

func (this *transactionService) GetByAccountId(ctx context.Context, userId string, accountId string) ([]entities.Transaction, error) {
	transactions, err := this.transactionRepository.GetByAccountId(ctx, accountId)
	if err != nil {
		this.logger.Errorf("Failed to get transactions: %v", err)
		return nil, err
	}

	return transactions, nil
}
