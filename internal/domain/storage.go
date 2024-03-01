package domain

import (
	"github.com/google/uuid"
	"time"
)

type Storage struct {
	Id          uuid.UUID
	Title       string
	Geolocation string
	TotalSize   float64
	Created     time.Time
}
