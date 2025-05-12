package services

import (
	"bank-system/src/repositories"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type SchedulerService interface {
	StartPaymentOverdueChecker(ctx context.Context)
	checkOverduePayments(ctx context.Context) error
}

type schedulerService struct {
	paymentRepository repositories.PaymentRepository
	logger            *logrus.Logger
}

func NewSchedulerService(paymentRepository repositories.PaymentRepository, logger *logrus.Logger) SchedulerService {
	return &schedulerService{
		paymentRepository: paymentRepository,
		logger:            logger,
	}
}

func (this *schedulerService) checkOverduePayments(ctx context.Context) error {
	this.logger.Info("Checking for overdue payments")

	// Get all payments that are due but not paid and not already marked as overdue
	payments, err := this.paymentRepository.GetOverdue(ctx)
	if err != nil {
		return err
	}

	this.logger.Infof("Found %d overdue payments", len(payments))

	// Mark each payment as overdue
	for _, payment := range payments {
		payment.Status = "OVERDUE"

		_, err = this.paymentRepository.Update(ctx, payment)

		if err != nil {
			this.logger.Errorf("Failed to mark payment %s as overdue: %v", payment.ID, err)
			continue
		}

		this.logger.Infof("Marked payment %s as overdue", payment.ID)
	}

	return nil
}

func (this *schedulerService) StartPaymentOverdueChecker(ctx context.Context) {
	this.logger.Info("Starting payment overdue checker scheduler")

	if err := this.checkOverduePayments(ctx); err != nil {
		this.logger.Errorf("Error checking overdue payments: %v", err)
	}

	ticker := time.NewTicker(12 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := this.checkOverduePayments(ctx); err != nil {
					this.logger.Errorf("Error checking overdue payments: %v", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				this.logger.Info("Payment overdue checker stopped")
				return
			}
		}
	}()
}
