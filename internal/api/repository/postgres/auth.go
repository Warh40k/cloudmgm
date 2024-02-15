package postgres

import (
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func (r *AuthPostgres) SignUp(user domain.User) (int, error) {
	panic("implement me")
}

func (r *AuthPostgres) SignIn(username, password string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *AuthPostgres) GenerateToken(username, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *AuthPostgres) CheckToken(token string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
