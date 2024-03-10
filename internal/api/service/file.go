package service

import (
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"io"
	"math"
	"mime/multipart"
	"os"
	"strings"
)

var (
	ErrTypeExceeded = errors.New("iter counts overflow")
)

const IterCount = math.MaxInt8

type FileService struct {
	repos repository.File
}

func (s FileService) UploadFile(volumePath string, file multipart.File, fileName string, fs afero.Fs) (string, error) {
	if fileName == "" {
		return "", errors.New("error empty filename")
	}
	// Get unique name for file
	exist, err := afero.Exists(fs, volumePath+"/"+fileName)
	if err != nil {
		return "", err
	}
	// If no filename existed before
	if exist {
		// try to alter name to get unique one
		dotIndex := strings.LastIndex(fileName, ".")
		var nameParts [2]string
		if dotIndex != -1 {
			nameParts = [2]string{fileName[:dotIndex], fileName[dotIndex:]}
		} else {
			nameParts = [2]string{fileName, ""}
		}
		var i int

		for i = 1; i < IterCount; i++ {
			fileName = fmt.Sprintf("%s(%d)%s", nameParts[0], i, nameParts[1])

			exist, err = afero.Exists(fs, volumePath+"/"+fileName)
			if err != nil {
				return "", err
			}
			if !exist {
				break
			}
		}
		if i == IterCount-1 {
			return "", ErrTypeExceeded
		}
	}

	dst, err := fs.Create(volumePath + "/" + fileName)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	return fileName, nil
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

	path := fileInfo.GetPath()
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
