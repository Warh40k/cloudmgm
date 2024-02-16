package domain

import (
	"github.com/google/uuid"
	"time"
)

type VirtualMachine struct {
	Id      uuid.UUID `json:"id,omitempty" db:"id"`
	Label   string    `json:"label,omitempty" db:"label"`
	OS      string    `json:"os,omitempty" db:"os"`
	Created time.Time `json:"date_create" db:"created"`
	Size    float64   `json:"size,omitempty" db:"size"`
}
