package services

import (
	"bank-system/src/entities"
	"bank-system/src/repositories"
	"bank-system/src/utils"
	"math/rand"
	"time"

	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type CardService interface {
	Create(ctx context.Context, userId string, data entities.CreateCardDto) (entities.CreateCardResponseDto, error)
	GetInfo(ctx context.Context, userId string, data entities.GetCardInfoDto) (entities.GetCardInfoResponseDto, error)
	Pay(ctx context.Context, userId string, data entities.PayCardDto) (string, error)
}

type cardService struct {
	accountService AccountService
	cardRepository repositories.CardRepository
	logger         *logrus.Logger
}

func NewCardService(
	accountService AccountService,
	cardRepository repositories.CardRepository,
	logger *logrus.Logger,
) CardService {
	return &cardService{
		accountService: accountService,
		cardRepository: cardRepository,
		logger:         logger,
	}
}

func (this *cardService) generateCardNumber() string {
	const cardNumberLength = 16
	cardNumber := make([]byte, cardNumberLength-1) // Generate 15 digits, last one will be checksum

	cardNumber[0] = byte(rand.Intn(3)+4) + '0'

	for i := 1; i < cardNumberLength-1; i++ {
		cardNumber[i] = byte(rand.Intn(10)) + '0'
	}

	sum := 0
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')
		if (len(cardNumber)-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	checkDigit := (10 - (sum % 10)) % 10

	return string(cardNumber) + string(byte(checkDigit)+'0')
}

func (this *cardService) generateCVV() string {
	const cvvLength = 3

	cvv := make([]byte, cvvLength)

	for i := range cvvLength {
		digit := rand.Intn(10)
		cvv[i] = byte(digit) + '0'
	}

	return string(cvv)
}

func (this *cardService) hashCVV(cvv string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (this *cardService) Create(ctx context.Context, userId string, data entities.CreateCardDto) (entities.CreateCardResponseDto, error) {
	account, err := this.accountService.GetById(ctx, data.AccountID)

	if err != nil {
		this.logger.Errorf("Failed to get account: %v", err)
		return entities.CreateCardResponseDto{}, err
	}
	if account.UserID != userId {
		this.logger.Error("Unauthorised")
		return entities.CreateCardResponseDto{}, errors.New("Unauthorised")
	}

	cardNumber := this.generateCardNumber()
	encryptedCardNumber, err := utils.EncryptWithPGP(
		cardNumber,
		data.PGPKey,
	)
	if err != nil {
		this.logger.Errorf("Failed to encrypt card number: %v", err)
		return entities.CreateCardResponseDto{}, err
	}

	expiration := time.Now().AddDate(5, 0, 0).Format("2006-01-02")
	encryptedCardExpiration, err := utils.EncryptWithPGP(expiration, data.PGPKey)
	if err != nil {
		this.logger.Errorf("Failed to encrypt card expiration: %v", err)
		return entities.CreateCardResponseDto{}, err
	}

	cvv := this.generateCVV()
	cvvHash, err := this.hashCVV(cvv)
	if err != nil {
		this.logger.Errorf("Failed to hash cvv: %v", err)
		return entities.CreateCardResponseDto{}, err
	}

	cardId, err := this.cardRepository.Create(
		ctx,
		entities.Card{
			ID:         uuid.New().String(),
			Number:     encryptedCardNumber,
			Expiration: encryptedCardExpiration,
			CVV:        cvvHash,
			UserId:     userId,
			AccountID:  data.AccountID,
			CreatedAt:  time.Now().Unix(),
		},
	)

	if err != nil {
		this.logger.Errorf("Failed to create card: %v", err)
		return entities.CreateCardResponseDto{}, err
	}

	createdCard := entities.CreateCardResponseDto{
		ID:         cardId,
		Number:     cardNumber,
		Expiration: expiration,
		CVV:        cvv,
	}
	return createdCard, nil
}

func (this *cardService) GetInfo(ctx context.Context, userId string, data entities.GetCardInfoDto) (entities.GetCardInfoResponseDto, error) {
	card, err := this.cardRepository.GetById(ctx, data.CardId)
	if err != nil {
		this.logger.Errorf("Failed to get card: %v", err)
		return entities.GetCardInfoResponseDto{}, err
	}

	if card.UserId != userId {
		this.logger.Error("Unauthorised")
		return entities.GetCardInfoResponseDto{}, errors.New("Unauthorised")
	}

	decryptedCardNumber, err := utils.DecryptWithPGP(card.Number, data.PGPKey)
	if err != nil {
		this.logger.Errorf("Failed to decrypt card number: %v", err)
		return entities.GetCardInfoResponseDto{}, err
	}
	decryptedCardExpiration, err := utils.DecryptWithPGP(card.Expiration, data.PGPKey)
	if err != nil {
		this.logger.Errorf("Failed to decrypt card expiration: %v", err)
		return entities.GetCardInfoResponseDto{}, err
	}

	decryptedCard := entities.GetCardInfoResponseDto{
		ID:         card.ID,
		Number:     decryptedCardNumber,
		Expiration: decryptedCardExpiration,
		UserId:     card.UserId,
		AccountID:  card.AccountID,
		CreatedAt:  card.CreatedAt,
	}

	return decryptedCard, nil
}

func (this *cardService) Pay(ctx context.Context, userId string, data entities.PayCardDto) (string, error) {
	card, err := this.cardRepository.GetById(ctx, data.CardId)
	if err != nil {
		this.logger.Errorf("Failed to get card: %v", err)
		return "", err
	}
	if card.UserId != userId {
		this.logger.Error("Unauthorised")
		return "", errors.New("Unauthorised")
	}

	decryptedCardExpiration, err := utils.DecryptWithPGP(card.Expiration, data.PGPKey)
	if err != nil {
		this.logger.Errorf("Failed to decrypt card expiration: %v", err)
		return "", err
	}

	expirationTime, err := time.Parse("2006-01-02", decryptedCardExpiration)
	if err != nil {
		this.logger.Errorf("Failed to parse card expiration: %v", err)
		return "", err
	}

	if time.Now().After(expirationTime) {
		this.logger.Error("Card expired")
		return "", errors.New("Card expired")
	}

	if decryptedCardExpiration != data.Expiration {
		this.logger.Error("Invalid card expiration")
		return "", errors.New("Invalid card expiration")
	}

	err = bcrypt.CompareHashAndPassword([]byte(card.CVV), []byte(data.CVV))
	if err != nil {
		this.logger.Errorf("Failed to compare cvv: %v", err)
		return "", errors.New("Invalid cvv")
	}

	account, err := this.accountService.GetById(ctx, card.AccountID)
	if err != nil {
		this.logger.Errorf("Failed to get account: %v", err)
		return "", err
	}

	if account.Balance < data.Amount {
		this.logger.Error("Insufficient funds")
		return "", errors.New("Insufficient funds")
	}

	transactionId, err := this.accountService.UpdateBalance(
		ctx,
		account.ID,
		userId,
		entities.UpdateAccountBalanceDto{
			Amount: data.Amount,
			Type:   "payment",
		},
	)

	return transactionId, err
}
