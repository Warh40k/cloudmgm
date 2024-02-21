package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
)

type VmService struct {
	repos repository.Vm
}

func (s VmService) CheckOwnership(userId, vmId uuid.UUID) error {
	return s.repos.CheckOwnership(userId, vmId)
}

func (s VmService) ListVm(userId uuid.UUID) ([]domain.VirtualMachine, error) {
	vms, err := s.repos.ListVm(userId)
	if err != nil {
		return nil, ErrInternal
	}
	return vms, nil
}

func (s VmService) GetVm(vmId uuid.UUID) (domain.VirtualMachine, error) {
	//TODO implement me
	panic("implement me")
}

func (s VmService) CreateVm(userId uuid.UUID, machine domain.VirtualMachine) error {
	return s.repos.CreateVm(userId, machine)
}

func (s VmService) DeleteVm(vmId uuid.UUID) error {
	return s.repos.DeleteVm(vmId)
}

func (s VmService) ModifyVm(vmId uuid.UUID, machine domain.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}

func NewVmService(repos repository.Vm) *VmService {
	return &VmService{repos: repos}
}
