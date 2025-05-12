package services

import (
	"bank-system/src/entities"
	"bank-system/src/repositories"
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type LoanService interface {
	Apply(ctx context.Context, userId string, data entities.LoanApplyDto) (entities.LoanResponseDto, error)
	GetSchedule(ctx context.Context, userId string, loanId string) ([]entities.Payment, error)
}

type loanService struct {
	loanRepository     repositories.LoanRepository
	paymentRepository  repositories.PaymentRepository
	accountService     AccountService
	transactionService TransactionService
	logger             *logrus.Logger
}

func NewLoanService(
	loanRepository repositories.LoanRepository,
	paymentRepository repositories.PaymentRepository,
	accountService AccountService,
	transactionService TransactionService,
	logger *logrus.Logger,
) LoanService {
	return &loanService{
		loanRepository:     loanRepository,
		paymentRepository:  paymentRepository,
		accountService:     accountService,
		transactionService: transactionService,
		logger:             logger,
	}
}

func round2(val float64) float64 {
	return math.Round(val*100) / 100
}

func toCents(val float64) int64 {
	return int64(math.Round(val * 100))
}

func calcMonthlyPayment(amount int64, term int, interestRate float64) float64 {
	amountFloat := float64(amount)
	termFloat := float64(term)
	monthlyRate := interestRate / 12.0 / 100.0

	if monthlyRate == 0 {
		return amountFloat / termFloat
	}

	r := monthlyRate
	n := termFloat

	pow := math.Pow(1+r, n)
	numerator := r * pow
	denominator := pow - 1

	if denominator == 0 {
		return 0
	}

	monthlyPayment := amountFloat * (numerator / denominator)

	return math.Round(monthlyPayment*100) / 100
}

func calcPaymentSchedule(
	rate float64,
	monthlyPayment float64,
	startDate time.Time,
	amount int64,
	term int,
	loanId string,
) []entities.Payment {
	schedule := make([]entities.Payment, 0, term)
	remainingPrincipal := float64(amount)
	monthlyRate := rate / 12.0 / 100.0

	for i := 0; i < term; i++ {
		dueDate := startDate.AddDate(0, i+1, 0)

		interestPart := round2(remainingPrincipal * monthlyRate)
		principalPart := round2(monthlyPayment - interestPart)

		if i == term-1 || (remainingPrincipal-principalPart) <= 0.009 {
			principalPart = round2(remainingPrincipal)
			monthlyPayment = round2(principalPart + interestPart)
		}

		payment := entities.Payment{
			ID:            uuid.New().String(),
			LoanId:        loanId,
			Amount:        toCents(monthlyPayment),
			DueDate:       dueDate.Unix(),
			PrincipalPart: toCents(principalPart),
			InterestPart:  toCents(interestPart),
			Status:        "new",
			IsPaid:        false,
		}
		schedule = append(schedule, payment)

		remainingPrincipal -= principalPart
		if remainingPrincipal <= 0 {
			break
		}
	}

	return schedule
}

func (this *loanService) Apply(ctx context.Context, userId string, data entities.LoanApplyDto) (entities.LoanResponseDto, error) {
	account, err := this.accountService.GetById(ctx, data.AccountId)
	if err != nil {
		return entities.LoanResponseDto{}, err
	}

	if account.UserID != userId {
		this.logger.Errorf("Unauthorised")
		return entities.LoanResponseDto{}, errors.New("Unauthorised")
	}

	keyRate, err := GetCentralBankRate()
	if err != nil {
		this.logger.Errorf("Failed to get key rate: %v", err)
		return entities.LoanResponseDto{}, err
	}

	interestRate := keyRate * 1.1
	monthlyPayment := calcMonthlyPayment(data.Amount, data.Term, interestRate)

	startDate := time.Now()
	loanId := uuid.New().String()

	payments := calcPaymentSchedule(
		interestRate,
		monthlyPayment,
		startDate,
		data.Amount,
		data.Term,
		loanId,
	)

	_, err = this.paymentRepository.CreateMany(ctx, payments)

	loan := entities.Loan{
		ID:           loanId,
		UserId:       userId,
		AccountId:    data.AccountId,
		Amount:       data.Amount,
		InterestRate: interestRate,
		Term:         data.Term,
		StartDate:    startDate.Unix(),
		Debt:         data.Amount,
	}

	_, err = this.loanRepository.Create(ctx, loan)
	if err != nil {
		this.logger.Errorf("Failed to create loan: %v", err)
		return entities.LoanResponseDto{}, err
	}
	_, err = this.paymentRepository.CreateMany(ctx, payments)
	if err != nil {
		this.logger.Errorf("Failed to create payments: %v", err)
		return entities.LoanResponseDto{}, err
	}

	this.accountService.UpdateBalance(
		ctx,
		data.AccountId,
		userId,
		entities.UpdateAccountBalanceDto{
			Amount: data.Amount,
			Type:   "loan",
		},
	)

	this.logger.Info()

	return entities.LoanResponseDto{
		ID:           loan.ID,
		UserId:       loan.UserId,
		AccountId:    loan.AccountId,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		Term:         loan.Term,
		StartDate:    loan.StartDate,
		Debt:         loan.Debt,
		Payments:     payments,
	}, nil
}

func (this *loanService) GetSchedule(ctx context.Context, userId string, loanId string) ([]entities.Payment, error) {
	loan, err := this.loanRepository.GetById(ctx, loanId)
	if err != nil {
		this.logger.Errorf("Failed to get loan: %v", err)
		return []entities.Payment{}, err
	}

	if loan.UserId != userId {
		this.logger.Errorf("Unauthorised")
		return []entities.Payment{}, errors.New("Unauthorised")
	}

	payments, err := this.paymentRepository.GetByLoanId(ctx, loanId)
	if err != nil {
		this.logger.Errorf("Failed to get payments: %v", err)
		return []entities.Payment{}, err
	}

	return payments, nil
}
