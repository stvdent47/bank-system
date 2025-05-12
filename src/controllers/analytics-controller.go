package controllers

import (
	"bank-system/src/entities"
	"bank-system/src/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type AnalyticsController struct {
	analyticsService services.AnalyticsService
	logger           *logrus.Logger
}

func NewAnalyticsController(
	analyticsService services.AnalyticsService,
	logger *logrus.Logger,
) *AnalyticsController {
	return &AnalyticsController{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

func (this *AnalyticsController) GetTransactionsAnalytics(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)
	var data entities.GetTransactionsAnalyticsDto

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		this.logger.Errorf("Failed to decode request body: %v", err)
		http.Error(
			w,
			fmt.Errorf("Failed to decode request body: %v", err).Error(),
			http.StatusBadRequest,
		)
		return
	}

	transactions, err := this.analyticsService.GetTransactionsAnalytics(r.Context(), userId, data.AccountId)
	if err != nil {
		this.logger.Errorf("Failed to get transactions: %v", err)
		http.Error(
			w,
			fmt.Errorf("Failed to get transactions: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}
