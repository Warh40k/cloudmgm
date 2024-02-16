package domain

import (
	"github.com/google/uuid"
	"time"
)

type VirtualMachine struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Label      string    `json:"label,omitempty"`
	OS         string    `json:"os,omitempty"`
	DateCreate time.Time `json:"date_create"`
	Size       float64   `json:"size,omitempty"`
}
