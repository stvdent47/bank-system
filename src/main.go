package main

import (
	"bank-system/src/controllers"
	"bank-system/src/db"
	"bank-system/src/middlewares"
	"bank-system/src/repositories"
	"bank-system/src/services"
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const PORT = "8888"

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	ctx := context.Background()

	db, err := db.New(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	// repositories
	userRepository := repositories.NewUserRepository(db)
	transactionRepository := repositories.NewTransactionRepository(db)
	accountRepository := repositories.NewAccountRepository(db)
	cardRepository := repositories.NewCardRepository(db)
	paymentRepository := repositories.NewPaymentRepository(db)
	loanRepository := repositories.NewLoanRepository(db)

	// services
	userService := services.NewUserService(
		userRepository,
		logger,
	)
	transactionService := services.NewTransactionService(
		transactionRepository,
		logger,
	)
	accountService := services.NewAccountService(
		accountRepository,
		transactionService,
		logger,
	)
	cardService := services.NewCardService(
		accountService,
		cardRepository,
		logger,
	)
	loanService := services.NewLoanService(
		loanRepository,
		paymentRepository,
		accountService,
		transactionService,
		logger,
	)
	analyticsService := services.NewAnalyticsService(
		accountService,
		transactionService,
		logger,
	)
	services.NewSchedulerService(paymentRepository, logger).StartPaymentOverdueChecker(ctx)

	// controllers
	userController := controllers.NewUserController(
		userService,
		logger,
	)
	accountController := controllers.NewAccountController(
		accountService,
		logger,
	)
	cardController := controllers.NewCardController(
		cardService,
		logger,
	)
	loanController := controllers.NewLoanController(
		loanService,
		logger,
	)
	analyticsController := controllers.NewAnalyticsController(
		analyticsService,
		logger,
	)

	jwtMiddleware := middlewares.NewJwtMiddleware(logger)

	// routes
	router := mux.NewRouter().PathPrefix("").Subrouter()
	// auth
	router.HandleFunc("/register", userController.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", userController.Login).Methods(http.MethodPost)
	// accounts
	accountRouter := router.PathPrefix("/accounts").Subrouter()
	accountRouter.Use(jwtMiddleware.Middleware)
	accountRouter.HandleFunc("", accountController.GetAll).Methods(http.MethodGet)
	accountRouter.HandleFunc("/{id}", accountController.GetById).Methods(http.MethodGet)
	accountRouter.HandleFunc("/create", accountController.Create).Methods(http.MethodPost)
	accountRouter.HandleFunc("/{id}/balance", accountController.UpdateBalance).Methods(http.MethodPatch)
	accountRouter.HandleFunc("/{id}", accountController.Delete).Methods(http.MethodDelete)
	accountRouter.HandleFunc("/transfer", accountController.Transfer).Methods(http.MethodPost)
	// cards
	cardRouter := router.PathPrefix("/cards").Subrouter()
	cardRouter.Use(jwtMiddleware.Middleware)
	cardRouter.HandleFunc("/create", cardController.Create).Methods(http.MethodPost)
	cardRouter.HandleFunc("/info", cardController.GetInfo).Methods(http.MethodPost)
	cardRouter.HandleFunc("/pay", cardController.Pay).Methods(http.MethodPost)
	// loans
	loanRouter := router.PathPrefix("/loans").Subrouter()
	loanRouter.Use(jwtMiddleware.Middleware)
	loanRouter.HandleFunc("/apply", loanController.Apply).Methods(http.MethodPost)
	loanRouter.HandleFunc("/{id}/schedule", loanController.GetSchedule).Methods(http.MethodGet)
	// analytics
	analyticsRouter := router.PathPrefix("/analytics").Subrouter()
	analyticsRouter.Use(jwtMiddleware.Middleware)
	analyticsRouter.HandleFunc("/transactions", analyticsController.GetTransactionsAnalytics).Methods(http.MethodPost)

	logger.Infof("Starting server at: localhost:%s", PORT)
	err = http.ListenAndServe(
		fmt.Sprintf(":%s", PORT),
		router,
	)

	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
		return
	}
}
