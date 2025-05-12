package controllers

import (
	"bank-system/src/entities"
	"bank-system/src/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type LoanController struct {
	loanService services.LoanService
	logger      *logrus.Logger
}

func NewLoanController(loanService services.LoanService, logger *logrus.Logger) *LoanController {
	return &LoanController{
		loanService: loanService,
		logger:      logger,
	}
}

func (this *LoanController) Apply(w http.ResponseWriter, r *http.Request) {
	var data entities.LoanApplyDto

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

	if !data.IsValid() {
		this.logger.Errorf("Invalid request body: %v", data)
		http.Error(
			w,
			fmt.Errorf("Invalid request body: %v", data).Error(),
			http.StatusBadRequest,
		)
		return
	}

	userId := r.Context().Value("userId").(string)
	loan, err := this.loanService.Apply(r.Context(), userId, data)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to apply loan: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}
	json.NewEncoder(w).Encode(loan)
}

func (this *LoanController) GetSchedule(w http.ResponseWriter, r *http.Request) {
	payments, err := this.loanService.GetSchedule(
		r.Context(),
		r.Context().Value("userId").(string),
		mux.Vars(r)["id"],
	)

	if err != nil {
		this.logger.Errorf("Failed to get loan schedule: %v", err)
		http.Error(
			w,
			fmt.Errorf("Failed to get loan schedule: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(payments)
}
