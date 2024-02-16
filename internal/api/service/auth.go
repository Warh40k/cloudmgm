package service

import (
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/api/api_errors"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/Warh40k/cloud-manager/internal/utils/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"time"
)

const (
	salt       = "fjlsj2374slfjsd728vvnts"
	tokenTTL   = 12 * time.Hour
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFsdsfdjH"
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (s AuthService) SignUp(user domain.User) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("error creating uuid: %w", err)
	}
	user.Id = id
	hash, err := auth.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	user.Password = hash

	err = s.repos.SignUp(user)
	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("login already exist %w", api_errors.ErrBadRequest)
			} else {
				return api_errors.ErrInternal
			}
		}
	}

	return nil
}

func (s AuthService) SignIn(username, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}
