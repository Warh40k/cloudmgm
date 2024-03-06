package service

import (
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"io"
	"math"
	"mime/multipart"
	"os"
	"strings"
)

var (
	ErrTypeExceeded = errors.New("iter counts overflow")
)

const IterCount = math.MaxInt16

type FileService struct {
	repos repository.File
}

func (s FileService) UploadFile(volumeId uuid.UUID, file *multipart.File, header *multipart.FileHeader) (string, error) {
	path := viper.GetString("files.save_path") + "/" + volumeId.String()

	// Get unique name for file
	_, err := os.Stat(path + "/" + header.Filename)
	if os.IsNotExist(err) {
		return header.Filename, nil
	}

	dotIndex := strings.LastIndex(header.Filename, ".")
	nameParts := []string{header.Filename[:dotIndex], header.Filename[dotIndex:]}
	var resultName string
	var i int

	for i = 1; i < IterCount; i++ {
		resultName = fmt.Sprintf("%s(%d)%s", nameParts[0], i, nameParts[1])
		if _, err = os.Stat(path + "/" + resultName); os.IsNotExist(err) {
			break
		}
	}

	if !(i < IterCount) {
		return "", ErrTypeExceeded
	}

	dst, err := os.Create(path + "/" + resultName)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, *file); err != nil {
		return "", err
	}

	return resultName, nil
}

func (s FileService) ListVolumeFiles(volumeId uuid.UUID) ([]domain.File, error) {
	return s.repos.ListVolumeFiles(volumeId)
}

func (s FileService) CreateFile(file domain.File) (uuid.UUID, error) {
	return s.repos.CreateFile(file)
}

func (s FileService) DeleteFile(fileId uuid.UUID) error {
	fileInfo, err := s.GetFileInfo(fileId)
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	path := viper.GetString("files.save_path") + "/" + fileInfo.VolumeId.String() + "/" + fileInfo.Name
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("specified file not found: %w", err)
	}
	if err = os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return s.repos.DeleteFileInfo(fileId)
}

func (s FileService) GetFileInfo(id uuid.UUID) (domain.File, error) {
	return s.repos.GetFile(id)
}

func (s FileService) SearchFile(filename string) ([]File, error) {
	panic("not implemented")
}

func NewFileService(repos repository.File) *FileService {
	return &FileService{repos: repos}
}
