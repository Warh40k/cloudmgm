package postgres

import (
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FilePostgres struct {
	db *sqlx.DB
}

func (r FilePostgres) ListVolumeFiles(volumeId uuid.UUID) ([]domain.File, error) {
	query := fmt.Sprintf(`SELECT * FROM %s where volume_id = $1`, filesTable)
	var files []domain.File
	err := r.db.Select(&files, query, volumeId)
	if err != nil {
		return nil, ErrInternal
	}
	return files, nil
}

func (r FilePostgres) CreateFile(file domain.File) (uuid.UUID, error) {
	file.Id = uuid.New()

	vmQuery := fmt.Sprintf(`INSERT INTO %s(id, volume_id, name, size, link) 
								VALUES($1,$2,$3,$4,$5)`, filesTable)
	_, err := r.db.Exec(vmQuery, file.Id, file.VolumeId, file.Name, file.Size, file.Link)
	if err != nil {
		return uuid.Nil, ErrInternal
	}
	return file.Id, nil
}

func (r FilePostgres) DeleteFile(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r FilePostgres) GetFile(fileId uuid.UUID) (domain.File, error) {
	var file domain.File
	query := fmt.Sprintf(`SELECT * FROM %s where id = $1`, filesTable)
	err := r.db.Get(&file, query, fileId)
	if err != nil {
		return file, ErrNoRows
	}
	return file, nil
}

func (r FilePostgres) SearchFile(filename string) ([]domain.File, error) {
	//TODO implement me
	panic("implement me")
}

func NewFilePostgres(db *sqlx.DB) *FilePostgres {
	return &FilePostgres{db: db}
}
