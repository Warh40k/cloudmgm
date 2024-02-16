package repository

import (
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *AuthPostgres) GetUserByLogin(login string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE login=$1`, usersTable)
	err := r.db.Get(&user, query, login)
	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) {
			return user, ErrInternal
		} else {
			return user, ErrNoRows
		}
	}
	return user, nil
}
