package services

import (
	"bank-system/src/entities"
	"bank-system/src/repositories"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AccountService interface {
	Create(ctx context.Context, userId string) (string, error)
	GetAll(ctx context.Context, userId string) ([]entities.Account, error)
	GetById(ctx context.Context, accountId string) (entities.Account, error)
	UpdateBalance(ctx context.Context, accountId string, userId string, data entities.UpdateAccountBalanceDto) (string, error)
	Delete(ctx context.Context, accountId string) (string, error)
	Transfer(ctx context.Context, userId string, data entities.TransferDto) (string, error)
}

type accountService struct {
	accountRepository  repositories.AccountRepository
	transactionService TransactionService
	logger             *logrus.Logger
}

func NewAccountService(
	accountRepository repositories.AccountRepository,
	transactionService TransactionService,
	logger *logrus.Logger,
) AccountService {
	return &accountService{
		accountRepository:  accountRepository,
		transactionService: transactionService,
		logger:             logger,
	}
}

func (this *accountService) Create(ctx context.Context, userId string) (string, error) {
	accountId, err := this.accountRepository.Create(
		ctx,
		entities.Account{
			ID:        uuid.New().String(),
			Balance:   0,
			UserID:    userId,
			CreatedAt: time.Now().Unix(),
		},
	)

	if err != nil {
		this.logger.Errorf("Failed to create account: %v", err)
		return "", err
	}

	return accountId, nil
}

func (this *accountService) GetAll(ctx context.Context, userId string) ([]entities.Account, error) {
	accounts, err := this.accountRepository.GetAll(ctx, userId)

	if err != nil {
		this.logger.Errorf("Failed to get accounts: %v", err)
		return nil, err
	}

	return accounts, nil
}

func (this *accountService) GetById(ctx context.Context, accountId string) (entities.Account, error) {
	account, err := this.accountRepository.GetById(ctx, accountId)

	if err != nil {
		this.logger.Errorf("Failed to get account: %v", err)
		return entities.Account{}, err
	}

	return account, nil
}

func (this *accountService) UpdateBalance(ctx context.Context, accountId string, userId string, data entities.UpdateAccountBalanceDto) (string, error) {
	account, err := this.accountRepository.GetById(ctx, accountId)

	if err != nil {
		this.logger.Errorf("Failed to get account: %v", err)
		return "", err
	}

	if account.UserID != userId {
		this.logger.Errorf("Unauthorised")
		return "", errors.New("Unauthorised")
	}

	var newBalance int64

	if data.Type == "deposit" {
		newBalance = account.Balance + data.Amount
	}
	if data.Type == "withdrawal" {
		newBalance = account.Balance - data.Amount
	}

	err = this.accountRepository.UpdateBalance(ctx, accountId, newBalance)
	if err != nil {
		this.logger.Errorf("Failed to update account: %v", err)
		return "", err
	}

	this.logger.Info("Balance updated: ", accountId)

	transactionId, err := this.transactionService.Create(
		ctx,
		userId,
		entities.Transaction{
			ID:          uuid.New().String(),
			Amount:      data.Amount,
			ToAccountId: accountId,
			Type:        data.Type,
			Description: "",
			CreatedAt:   time.Now().Unix(),
		},
	)
	if err != nil {
		this.logger.Errorf("Failed to create transaction: %v", err)
		return "", err
	}

	this.logger.Info("Transaction created: ", transactionId)

	return transactionId, nil
}

func (this *accountService) Delete(ctx context.Context, accountId string) (string, error) {
	deletedAccountId, err := this.accountRepository.Delete(ctx, accountId)

	if err != nil {
		this.logger.Errorf("Failed to delete account: %v", err)
		return "", err
	}

	return deletedAccountId, nil
}

func (this *accountService) Transfer(ctx context.Context, userId string, data entities.TransferDto) (string, error) {
	account, err := this.accountRepository.GetById(ctx, data.FromAccountId)
	if err != nil {
		this.logger.Errorf("Failed to get account: %v", err)
		return "", err
	}
	if account.UserID != userId {
		this.logger.Errorf("Unauthorised")
		return "", errors.New("Unauthorised")
	}
	if account.Balance < data.Amount {
		this.logger.Errorf("Insufficient balance")
		return "", errors.New("Insufficient balance")
	}

	err = this.accountRepository.Transfer(ctx, data.FromAccountId, data.ToAccountId, data.Amount)
	if err != nil {
		this.logger.Errorf("Failed to transfer money: %v", err)
		return "", err
	}

	this.logger.Info("Transfer successful")

	transactionId, err := this.transactionService.Create(
		ctx,
		userId,
		entities.Transaction{
			ID:            uuid.New().String(),
			Amount:        data.Amount,
			FromAccountId: data.FromAccountId,
			ToAccountId:   data.ToAccountId,
			Type:          "transfer",
			Description:   data.Description,
			CreatedAt:     time.Now().Unix(),
		},
	)

	return transactionId, nil
}
