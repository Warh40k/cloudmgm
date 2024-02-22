package domain

import (
	"github.com/google/uuid"
	"time"
)

type VirtualMachine struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Label       string    `json:"title" db:"title" validate:"required"`
	Description string    `json:"description" db:"description" validate:"required"`
	Created     time.Time `json:"created" db:"created"`
	Status      int       `json:"status" db:"status"`
	Size        float64   `json:"size" db:"size"`
}
