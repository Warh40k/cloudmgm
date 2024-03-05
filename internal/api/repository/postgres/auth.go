package postgres

import (
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/api/repository/response"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
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

func (r *AuthPostgres) SignUp(user domain.User) (uuid.UUID, error) {
	query := fmt.Sprintf(`INSERT INTO %s(id,name,username,password_hash) VALUES($1,$2,$3,$4) RETURNING id`, usersTable)
	row := r.db.QueryRowx(query, user.Id, user.Name, user.Username, user.PasswordHash)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueErrCode {
			return uuid.Nil, response.ErrUnique
		}

		return uuid.Nil, response.ErrInternal
	}
	return id, nil
}

func (r *AuthPostgres) GetUserByUsername(username string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE username=$1`, usersTable)
	err := r.db.Get(&user, query, username)
	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) {
			return user, response.ErrInternal
		} else {
			return user, response.ErrNoRows
		}
	}
	return user, nil
}
