package postgres

import (
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FilePostgres struct {
	db *sqlx.DB
}

func (r FilePostgres) CreateFile(file domain.File) error {
	//TODO implement me
	panic("implement me")
}

func (r FilePostgres) DeleteById(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r FilePostgres) GetById(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r FilePostgres) SearchFile(filename string) ([]domain.File, error) {
	//TODO implement me
	panic("implement me")
}

func NewFilePostgres(db *sqlx.DB) *FilePostgres {
	return &FilePostgres{db: db}
}
