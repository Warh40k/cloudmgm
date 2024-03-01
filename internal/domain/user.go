package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id           uuid.UUID `json:"id,omitempty" db:"id"`
	Name         string    `json:"name" db:"name" validate:"required"`
	Username     string    `json:"username" db:"username" validate:"required"`
	Password     string    `json:"password" db:"-" validate:"required"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Balance      float64   `json:"balance" db:"balance" validate:"gte=0"`
	Created      time.Time `json:"created" db:"created"`
}
