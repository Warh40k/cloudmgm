package postgres

import (
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	uniqueErrorCode = "23505"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func (r *AuthPostgres) SignUp(user domain.User) error {
	query := fmt.Sprintf(`INSERT INTO %s(id,name,login,password_hash) VALUES($1,$2,$3,$4)`, usersTable)
	_, err := r.db.Exec(query, user.Id, user.Name, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

func (r *AuthPostgres) SignIn(username, password string) (string, error) {
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
