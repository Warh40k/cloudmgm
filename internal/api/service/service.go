package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
)

import (
	"github.com/google/uuid"
)

type Service struct {
	Authorization
	Volume
	File
}

type Authorization interface {
	SignUp(user domain.User) error
	SignIn(username, password string) (string, error)
}

type Volume interface {
	ListVolume(userId uuid.UUID) ([]domain.Volume, error)
	GetVolume(vmId uuid.UUID) (domain.Volume, error)
	CreateVolume(userId uuid.UUID, machine domain.Volume) (uuid.UUID, error)
	DeleteVolume(vmId uuid.UUID) error
	UpdateVolume(machine domain.Volume) error
	CheckOwnership(userId, vmId uuid.UUID) error
	ResizeVolume(userId uuid.UUID, amount int) error
}

type File interface {
	CreateFile(file domain.File) (uuid.UUID, error)
	DeleteById(id uuid.UUID) error
	GetById(id uuid.UUID) error
	SearchFile(filename string) ([]File, error)
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Volume:        NewVolumeService(repos.Volume),
		File:          NewFileService(repos.File),
	}
}
