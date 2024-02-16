package repository

import (
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	SignUp(user domain.User) error
	GetUserByLogin(login string) (domain.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
