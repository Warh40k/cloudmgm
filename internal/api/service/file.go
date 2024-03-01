package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
)

type FileService struct {
	repos repository.File
}

func (s FileService) CreateFile(file domain.File) (uuid.UUID, error) {
	return s.repos.CreateFile(file)
}

func (s FileService) DeleteById(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s FileService) GetById(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s FileService) SearchFile(filename string) ([]File, error) {
	//TODO implement me
	panic("implement me")
}

func NewFileService(repos repository.File) *FileService {
	return &FileService{repos: repos}
}
