package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
)

type Service struct {
}

type Authorization interface {
	Register(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	CheckToken(token string) (bool, error)
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
