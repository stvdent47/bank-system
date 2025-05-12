package utils

import (
	"bank-system/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Duration(config.LoadJwtConfig().Exp) * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	jwt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.LoadJwtConfig().Secret))

	if err != nil {
		return "", err
	}

	return jwt, nil
}

func ValidateJwt(jwtString string) (string, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(config.LoadJwtConfig().Secret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("Invalid token claims")
	}

	userId, ok := claims["sub"].(string)

	if !ok {
		return "", errors.New("Invalid user id")
	}

	return userId, nil
}
