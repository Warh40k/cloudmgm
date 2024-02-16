package repository

import (
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

const (
	uniqueErrCode = "23505"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) SignUp(user domain.User) error {
	query := fmt.Sprintf(`INSERT INTO %s(id,name,login,password_hash) VALUES($1,$2,$3,$4)`, usersTable)
	_, err := r.db.Exec(query, user.Id, user.Name, user.Login, user.PasswordHash)
	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueErrCode {
				return ErrUnique
			} else {
				return ErrInternal
			}
		}
	}
	return nil
}
