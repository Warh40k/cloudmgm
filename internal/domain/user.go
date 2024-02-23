package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id           uuid.UUID `json:"id,omitempty" db:"id"`
	Name         string    `json:"name" db:"name" validate:"required"`
	Login        string    `json:"login" validate:"required"`
	Password     string    `json:"password" db:"-" validate:"required"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Created      time.Time `json:"created" db:"created"`
}
