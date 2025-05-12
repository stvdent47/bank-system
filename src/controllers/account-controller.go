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

type AccountController struct {
	accountService services.AccountService
	logger         *logrus.Logger
}

func NewAccountController(accountService services.AccountService, logger *logrus.Logger) *AccountController {
	return &AccountController{
		accountService: accountService,
		logger:         logger,
	}
}

func (this *AccountController) Create(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)

	accountId, err := this.accountService.Create(r.Context(), userId)

	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to create account: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"accountId": accountId,
	})
}

func (this *AccountController) GetAll(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)

	accounts, err := this.accountService.GetAll(r.Context(), userId)

	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to get accounts: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(accounts)
}

func (this *AccountController) GetById(w http.ResponseWriter, r *http.Request) {
	accountId, err := this.accountService.GetById(r.Context(), mux.Vars(r)["id"])

	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to get account: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(accountId)
}

func (this *AccountController) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	accountId := mux.Vars(r)["id"]
	userId := r.Context().Value("userId").(string)
	var data entities.UpdateAccountBalanceDto

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		this.logger.Errorf("Failed to decode request body: %v", err)
		http.Error(
			w,
			fmt.Errorf("Failed to parse request body: %v", err).Error(),
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

	transactionId, err := this.accountService.UpdateBalance(
		r.Context(),
		accountId,
		userId,
		data,
	)

	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to update balance: %v", err).Error(),
			http.StatusInternalServerError,
		)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"transactionId": transactionId,
	})
}

func (this *AccountController) Delete(w http.ResponseWriter, r *http.Request) {
	accountId, err := this.accountService.Delete(r.Context(), mux.Vars(r)["id"])

	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to delete account: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"accountId": accountId,
	})
}

func (this *AccountController) Transfer(w http.ResponseWriter, r *http.Request) {
	var data entities.TransferDto

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		this.logger.Errorf("Failed to decode request body: %v", err)
		http.Error(
			w,
			fmt.Errorf("Failed to parse request body: %v", err).Error(),
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

	transactionId, err := this.accountService.Transfer(r.Context(), userId, data)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to transfer money: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"transactionId": transactionId,
	})
}
