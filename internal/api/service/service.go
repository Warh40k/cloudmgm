package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"log/slog"
	"mime/multipart"
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
	DeleteFile(fileId uuid.UUID) error
	GetFileInfo(id uuid.UUID) (domain.File, error)
	ListVolumeFiles(volumeId uuid.UUID) ([]domain.File, error)
	SearchFile(filename string) ([]File, error)
	UploadFile(volumeId uuid.UUID, file *multipart.File, header *multipart.FileHeader) (string, error)
}

func NewService(repos *repository.Repository, log *slog.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Volume:        NewVolumeService(repos.Volume),
		File:          NewFileService(repos.File),
	}
}
