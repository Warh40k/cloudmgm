package service

import (
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/api/service/utils"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (s *AuthService) Pong() string {
	return "Pong"
}

func (s *AuthService) SignUp(user domain.User) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("error creating uuid: %w", err)
	}
	user.Id = id
	hash, err := utils.HashPassword(user.PasswordHash)
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

func (s *AuthService) SignIn(login, password string) (string, error) {
	user, err := s.repos.GetUserByLogin(login)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return "", ErrNotFound
		} else {
			return "", ErrInternal
		}
	}
	hash := user.PasswordHash
	if utils.CheckPassword(password, hash) {
		token, err := utils.GenerateJWT(user)
		if err != nil {
			return "", ErrInternal
		}
		return token, nil
	}

	return "", ErrUnauthorized
}
