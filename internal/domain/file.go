package domain

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	Id       uuid.UUID
	VolumeId uuid.UUID
	Name     string
	Link     string
	Created  time.Time
}
