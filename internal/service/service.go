package service

import (
	"github.com/Warh40k/cloud-manager/internal/repository"
)

type Service struct {
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
