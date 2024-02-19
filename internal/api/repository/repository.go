package repository

import (
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
	GetVm(userId uuid.UUID) (domain.VirtualMachine, error)
	CreateVm(userId uuid.UUID, machine domain.VirtualMachine) error
	DeleteVm(userId uuid.UUID) error
	ModifyVm(userId uuid.UUID, machine domain.VirtualMachine) error
}

type Repository struct {
	Authorization
	Vm
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Vm:            NewVmPostgres(db),
	}
}
