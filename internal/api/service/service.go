package service

import (
	"errors"
	"github.com/Warh40k/cloud-manager/internal/api/cache"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/spf13/afero"
	"log/slog"
	"mime/multipart"
)

import (
	"github.com/google/uuid"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrInternal     = errors.New("internal errors")
	ErrUnauthorized = errors.New("not authorized")
)

//go:generate mockery --all --dry-run=false
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
	SearchFile(filename string) ([]domain.File, error)
	UploadFile(volumePath string, file multipart.File, fileName string, fs afero.Fs) (string, error)
	ComposeZipArchive(files []domain.File, fs afero.Fs) (string, error)
	//GetFileInfo(fileId uuid.UUID) (multipart.File, error)
}

func NewService(repos *repository.Repository, cache *cache.Cache, log *slog.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, cache, log),
		Volume:        NewVolumeService(repos.Volume, cache, log),
		File:          NewFileService(repos.File, cache, log),
	}
}
