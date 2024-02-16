package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id           uuid.UUID `json:"id,omitempty" db:"id"`
	Name         string    `json:"name" db:"name"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"password" db:"password_hash"`
	Created      time.Time `json:"created" db:"created"`
}
