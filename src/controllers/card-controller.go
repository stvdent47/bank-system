package controllers

import (
	"bank-system/src/entities"
	"bank-system/src/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type CardController struct {
	cardService services.CardService
	logger      *logrus.Logger
}

func NewCardController(cardService services.CardService, logger *logrus.Logger) *CardController {
	return &CardController{
		cardService: cardService,
		logger:      logger,
	}
}

func (this *CardController) Create(w http.ResponseWriter, r *http.Request) {
	var data entities.CreateCardDto

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

	card, err := this.cardService.Create(r.Context(), userId, data)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to create card: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(card)
}

func (this *CardController) GetInfo(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)

	var data entities.GetCardInfoDto
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

	card, err := this.cardService.GetInfo(r.Context(), userId, data)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to get card info: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(card)
}

func (this *CardController) Pay(w http.ResponseWriter, r *http.Request) {
	var data entities.PayCardDto

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

	transactionId, err := this.cardService.Pay(r.Context(), userId, data)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to pay: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"transactionId": transactionId,
	})
}
