package repository

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository/postgres"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	SignUp(user domain.User) (uuid.UUID, error)
	GetUserByUsername(username string) (domain.User, error)
}

type Volume interface {
	ListVolume(userId uuid.UUID) ([]domain.Volume, error)
	GetVolume(vmId uuid.UUID) (domain.Volume, error)
	CreateVolume(userId uuid.UUID, machine domain.Volume) (uuid.UUID, error)
	DeleteVolume(vmId uuid.UUID) error
	UpdateVolume(machine domain.Volume) error
	CheckOwnership(userId, vmId uuid.UUID) error
}

type File interface {
	CreateFile(file domain.File) (uuid.UUID, error)
	DeleteById(id uuid.UUID) error
	GetById(id uuid.UUID) error
	SearchFile(filename string) ([]domain.File, error)
}

type Repository struct {
	Authorization
	Volume
	File
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(db),
		Volume:        postgres.NewVolumePostgres(db),
		File:          postgres.NewFilePostgres(db),
	}
}
