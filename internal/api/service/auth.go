package service

import (
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/api/service/utils"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type AuthService struct {
	repos repository.Authorization
	log   *slog.Logger
}

func NewAuthService(repos repository.Authorization, log *slog.Logger) *AuthService {
	return &AuthService{repos: repos, log: log}
}

func (s *AuthService) SignUp(user domain.User) error {
	if len(user.Password) < 8 {
		return ErrBadRequest
	}
	id := uuid.New()
	user.Id = id
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	user.PasswordHash = hash

	_, err = s.repos.SignUp(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SignIn(username, password string) (string, error) {
	user, err := s.repos.GetUserByUsername(username)
	if err != nil {
		return "", err
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
