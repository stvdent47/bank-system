package services

import (
	"bank-system/src/entities"
	"bank-system/src/repositories"
	"bank-system/src/utils"
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, data entities.RegisterUserDto) (string, error)
	Login(ctx context.Context, data entities.LoginUserDto) (string, error)
}

type userService struct {
	userRepository repositories.UserRepository
	logger         *logrus.Logger
}

func NewUserService(userRepository repositories.UserRepository, logger *logrus.Logger) UserService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (this *userService) Register(ctx context.Context, data entities.RegisterUserDto) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if err != nil {
		this.logger.Errorf("Failed to hash password: %v", err)
		return "", err
	}

	userId, err := this.userRepository.Create(
		ctx,
		&entities.User{
			ID:       uuid.New().String(),
			Username: data.Username,
			Password: string(passwordHash),
			Email:    data.Email,
		},
	)

	if err != nil {
		this.logger.Errorf("Failed to create user: %v", err)
		return "", err
	}

	go func() {
		err = SendEmailNotification(
			data.Email,
			"Registration",
			"You have successfully registered",
		)
		if err != nil {
			this.logger.Errorf("Failed to send email notification: %v", err)
		}
	}()

	return userId, nil
}

func (this *userService) Login(ctx context.Context, data entities.LoginUserDto) (string, error) {
	user, err := this.userRepository.FindByUsername(ctx, data.Username)

	if err != nil {
		this.logger.Errorf("Failed to find user: %v", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))

	if err != nil {
		this.logger.Errorf("Failed to compare password: %v", err)
		return "", err
	}

	jwt, err := utils.GenerateJwt(user.ID)
	if err != nil {
		this.logger.Errorf("Failed to generate jwt: %v", err)
		return "", err
	}

	return jwt, nil
}
