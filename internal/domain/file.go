package domain

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	Id       uuid.UUID
	VolumeId uuid.UUID
	Name     string
	Size     int64
	Link     string
	Created  time.Time
}
