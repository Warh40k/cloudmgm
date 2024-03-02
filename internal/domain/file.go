package domain

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	Id       uuid.UUID `json:"id" db:"id"`
	VolumeId uuid.UUID `json:"volume_id" db:"volume_id"`
	Name     string    `json:"name" db:"name"`
	Size     int64     `json:"size" db:"size"`
	Link     string    `json:"link" db:"link"`
	Created  time.Time `json:"created" db:"created"`
}
