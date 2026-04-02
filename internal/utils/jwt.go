package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var defaultSecret = []byte("supersecretkey_change_in_prod")

// GenerateToken gera um novo JWT para um usuário específico
func GenerateToken(userID string, secret []byte) (string, error) {
	if len(secret) == 0 {
		secret = defaultSecret
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(secret)
}

// ValidateToken valida o JWT e retorna o ID do usuário
func ValidateToken(tokenString string, secret []byte) (string, error) {
	if len(secret) == 0 {
		secret = defaultSecret
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metodo de assinatura invalido")
		}
		return secret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("user_id invalido no token")
		}
		return userID, nil
	}

	return "", errors.New("token invalido")
}
