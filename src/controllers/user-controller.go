package controllers

import (
	"bank-system/src/entities"
	"bank-system/src/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type UserController struct {
	userService services.UserService
	logger      *logrus.Logger
}

func NewUserController(userService services.UserService, logger *logrus.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

func (this *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var data entities.RegisterUserDto

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		this.logger.Errorf("failed to decode request body: %v", err)
		http.Error(
			w,
			fmt.Errorf("failed to decode request body: %v", err).Error(),
			http.StatusBadRequest,
		)
		return
	}

	if !data.IsValid() {
		http.Error(
			w,
			fmt.Errorf("username, password and email are required").Error(),
			http.StatusBadRequest,
		)
		return
	}

	userId, err := this.userService.Register(r.Context(), data)

	if err != nil {
		http.Error(
			w,
			fmt.Errorf("failed to create user: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"id": userId,
	})
}

func (this *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var data entities.LoginUserDto

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
		http.Error(
			w,
			fmt.Errorf("Username and password are required").Error(),
			http.StatusBadRequest,
		)
		return
	}

	jwt, err := this.userService.Login(r.Context(), data)

	if err != nil {
		http.Error(
			w,
			fmt.Errorf("Failed to login: %v", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"jwt": jwt,
	})
}
