package domain

import "time"

type VirtualMachine struct {
	Label      string    `json:"label,omitempty"`
	OS         string    `json:"os,omitempty"`
	DateCreate time.Time `json:"date_create"`
	Size       float64   `json:"size,omitempty"`
}
