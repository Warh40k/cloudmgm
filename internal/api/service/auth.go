package service

import (
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	salt       = "fjlsj2374slfjsd728vvnts"
	tokenTTL   = 12 * time.Hour
	signingKey = "j370sdfs34472fshvlruso043275fhka"
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
}

func GenerateJWT(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
	}, user.Id})
	return token.SignedString([]byte(signingKey))
}

func CheckJWT(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s AuthService) SignUp(user domain.User) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("error creating uuid: %w", err)
	}
	user.Id = id
	hash, err := HashPassword(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	user.PasswordHash = hash

	err = s.repos.SignUp(user)
	if err != nil {
		if errors.Is(err, repository.ErrUnique) {
			return ErrBadRequest
		} else {
			return ErrInternal
		}
	}

	return nil
}

func (s AuthService) SignIn(login, password string) (string, error) {
	user, err := s.repos.GetUserByLogin(login)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return "", ErrNotFound
		} else {
			return "", ErrInternal
		}
	}
	hash := user.PasswordHash
	if CheckPassword(password, hash) {
		token, err := GenerateJWT(user)
		if err != nil {
			return "", ErrInternal
		}
		return token, nil
	}

	return "", ErrUnauthorized
}
