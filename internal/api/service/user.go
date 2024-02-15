package service

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	salt       = "fjlsj2374slfjsd728vvnts"
	tokenTTL   = 12 * time.Hour
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFsdsfdjH"
)

func GenerateJWT() (string, error) {
	return "", nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
