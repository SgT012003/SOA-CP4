package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword gera o hash bcrypt para uma senha em texto plano.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash compara uma senha em texto plano com um hash bcrypt.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
