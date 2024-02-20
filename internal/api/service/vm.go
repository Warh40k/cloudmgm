package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
)

type VmService struct {
	repos repository.Vm
}

func (v VmService) ListVm(userId uuid.UUID) ([]domain.VirtualMachine, error) {
	vms, err := v.repos.ListVm(userId)
	if err != nil {
		return nil, ErrInternal
	}
	return vms, nil
}

func (v VmService) GetVm(userId uuid.UUID) (domain.VirtualMachine, error) {
	//TODO implement me
	panic("implement me")
}

func (v VmService) CreateVm(userId uuid.UUID, machine domain.VirtualMachine) error {
	return v.repos.CreateVm(userId, machine)
}

func (v VmService) DeleteVm(userId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (v VmService) ModifyVm(id uuid.UUID, machine domain.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}

func NewVmService(repos repository.Vm) *VmService {
	return &VmService{repos: repos}
}
