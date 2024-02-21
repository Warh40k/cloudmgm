package repository

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository/postgres"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	SignUp(user domain.User) error
	GetUserByLogin(login string) (domain.User, error)
}

type Vm interface {
	ListVm(userId uuid.UUID) ([]domain.VirtualMachine, error)
	GetVm(vmId uuid.UUID) (domain.VirtualMachine, error)
	CreateVm(userId uuid.UUID, machine domain.VirtualMachine) error
	DeleteVm(vmId uuid.UUID) error
	ModifyVm(vmId uuid.UUID, machine domain.VirtualMachine) error
	CheckOwnership(userId, vmId uuid.UUID) error
}

type Repository struct {
	Authorization
	Vm
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(db),
		Vm:            postgres.NewVmPostgres(db),
	}
}
