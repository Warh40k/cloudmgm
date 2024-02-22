package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
)

type Service struct {
	Authorization
	Vm
}

type Authorization interface {
	SignUp(user domain.User) error
	SignIn(username, password string) (string, error)
	Pong() string
}

type Vm interface {
	ListVm(userId uuid.UUID) ([]domain.VirtualMachine, error)
	GetVm(vmId uuid.UUID) (domain.VirtualMachine, error)
	CreateVm(userId uuid.UUID, machine domain.VirtualMachine) error
	DeleteVm(vmId uuid.UUID) error
	ModifyVm(vmId uuid.UUID, machine domain.VirtualMachine) error
	CheckOwnership(userId, vmId uuid.UUID) error
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Vm:            NewVmService(repos.Vm),
	}
}
