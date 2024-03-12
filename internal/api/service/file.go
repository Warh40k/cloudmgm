package service

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"io"
	"log/slog"
	"math"
	"mime/multipart"
	"os"
	"strings"
	"sync"
)

var (
	ErrTypeExceeded = errors.New("iter counts overflow")
)

const IterCount = math.MaxInt8

type FileService struct {
	repos repository.File
	log   *slog.Logger
}

func NewFileService(repos repository.File, log *slog.Logger) *FileService {
	return &FileService{repos: repos, log: log}
}

type fileContainer struct {
	file   *afero.File
	header domain.File
}

func fileConsumer(writer *zip.Writer, data <-chan fileContainer, done chan<- bool, errs chan<- error) {
	for f := range data {
		zipFile, err := writer.Create(f.header.Filename)
		if err != nil {
			errs <- fmt.Errorf("failed to create file in archive")
			done <- false
			return
		}
		_, err = io.Copy(zipFile, *f.file)
		if err != nil {
			errs <- fmt.Errorf("failed to write file in archive")
			done <- false
			return
		}
		(*f.file).Close()
	}
	done <- true
}

func (s FileService) ComposeZipArchive(fileHeaders []domain.File, fs afero.Fs) (string, error) {
	// нужна сервисная горутина, которая будет получать данные из канала и последовательно записывать их в архив
	// создать подкаталог в fs.TempDir и там архив
	// archive, err := afero.TempFile(fs, "", fmt.Sprintf("%s.zip", time.Now().String()))
	const op = "File.Service.ComposeZipArchive"
	log := s.log.With(
		slog.String("op", op),
	)

	//archivePath := strings.Join([]string{afero.GetTempDir(fs, ""), "archive.zip"}, "/")
	archive, err := afero.TempFile(fs, "", "*.zip")
	if err != nil {
		return "", err
	}

	zipWriter := zip.NewWriter(archive)
	defer func() {
		zipWriter.Close()
		archive.Close()
	}()

	data := make(chan fileContainer, len(fileHeaders))
	done := make(chan bool)
	consumerErrs := make(chan error, len(fileHeaders))
	senderErrs := make(chan error, len(fileHeaders))
	wg := sync.WaitGroup{}

	go fileConsumer(zipWriter, data, done, consumerErrs)

	for _, header := range fileHeaders {
		wg.Add(1)
		go func(header domain.File, wg *sync.WaitGroup) {
			defer wg.Done()
			file, err := fs.Open(header.GetPath())
			if err != nil {
				senderErrs <- fmt.Errorf("failed to open file %s", header.Filename)
				return
			}
			data <- struct {
				file   *afero.File
				header domain.File
			}{file: &file, header: header}
		}(header, &wg)
	}

	wg.Wait()
	close(data)
	close(senderErrs)
	ok := <-done
	close(consumerErrs)

	for err = range senderErrs {
		log.With(slog.String("err", err.Error())).Error("failed to send file to zip")
	}

	for err = range consumerErrs {
		log.With(slog.String("err", err.Error())).Error("failed to write file in zip")
	}
	if !ok {
		err = fmt.Errorf("failed to create archive")
	}
	return archive.Name(), err
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

func (s FileService) SearchFile(filename string) ([]domain.File, error) {
	panic("not implemented")
}
