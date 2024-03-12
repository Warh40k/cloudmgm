package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type VolumeService struct {
	repos repository.Volume
	log   *slog.Logger
}

func NewVolumeService(repos repository.Volume, log *slog.Logger) *VolumeService {
	return &VolumeService{repos: repos, log: log}
}

func (s VolumeService) ResizeVolume(userId uuid.UUID, amount int) error {
	//TODO implement me
	panic("implement me")
}

func (s VolumeService) CheckOwnership(userId, vmId uuid.UUID) error {
	return s.repos.CheckOwnership(userId, vmId)
}

func (s VolumeService) ListVolume(userId uuid.UUID) ([]domain.Volume, error) {
	vms, err := s.repos.ListVolume(userId)
	if err != nil {
		return nil, ErrInternal
	}
	return vms, nil
}

func (s VolumeService) GetVolume(vmId uuid.UUID) (domain.Volume, error) {
	return s.repos.GetVolume(vmId)
}

func (s VolumeService) CreateVolume(userId uuid.UUID, machine domain.Volume) (uuid.UUID, error) {
	return s.repos.CreateVolume(userId, machine)
}

func (s VolumeService) DeleteVolume(vmId uuid.UUID) error {
	return s.repos.DeleteVolume(vmId)
}

func (s VolumeService) UpdateVolume(machine domain.Volume) error {
	return s.repos.UpdateVolume(machine)
}
