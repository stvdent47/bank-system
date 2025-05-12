package services

import (
	"bank-system/src/entities"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type AnalyticsService interface {
	GetTransactionsAnalytics(ctx context.Context, userId string, accountId string) ([]entities.Transaction, error)
}

type analyticsService struct {
	accountService     AccountService
	transactionService TransactionService
	logger             *logrus.Logger
}

func NewAnalyticsService(
	accountService AccountService,
	transactionService TransactionService,
	logger *logrus.Logger,
) AnalyticsService {
	return &analyticsService{
		accountService:     accountService,
		transactionService: transactionService,
		logger:             logger,
	}
}

func (this *analyticsService) GetTransactionsAnalytics(ctx context.Context, userId string, accountId string) ([]entities.Transaction, error) {
	account, err := this.accountService.GetById(ctx, accountId)
	if err != nil {
		this.logger.Errorf("Failed to get account: %v", err)
		return nil, err
	}

	if account.UserID != userId {
		this.logger.Errorf("Unauthorised")
		return nil, errors.New("Unauthorised")
	}

	transactions, err := this.transactionService.GetByAccountId(ctx, userId, accountId)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
