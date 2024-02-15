package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
)

type Service struct {
	Authorization
}

type Authorization interface {
	SignUp(user domain.User) (int, error)
	SignIn(username, password string) (int, error)
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
