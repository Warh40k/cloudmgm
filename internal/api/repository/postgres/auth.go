package postgres

import (
	"errors"
	"fmt"
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
	query := fmt.Sprintf(`INSERT INTO %s(id,name,login,password_hash) VALUES($1,$2,$3,$4) RETURNING id`, usersTable)
	row := r.db.QueryRowx(query, user.Id, user.Name, user.Login, user.PasswordHash)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueErrCode {
			return uuid.Nil, ErrUnique
		}

		return uuid.Nil, ErrInternal
	}
	return id, nil
}
