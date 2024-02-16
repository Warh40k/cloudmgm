package repository

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository/postgres"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	SignUp(user domain.User) error
	SignIn(username, password string) (string, error)
	//GenerateToken(username, password string) (string, error)
	CheckToken(token string) (bool, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(db),
	}
}
