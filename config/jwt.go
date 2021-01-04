package config

import (
	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

var TokenAuth *jwtauth.JWTAuth

const (
	SecretKey string = "SECRET_KEY"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckHashPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
