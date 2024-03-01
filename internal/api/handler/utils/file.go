package utils

import (
	"errors"
	"fmt"
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

// UploadFile Create path directories and copy file to destination path
func UploadFile(storagePath string, file multipart.File, dst *os.File) error {
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}

// GetFileName Get unique filename for file (via num suffics)
func GetFileName(path, baseFileName string) (string, error) {
	// if no file with the same name in storage
	_, err := os.Stat(path + "/" + baseFileName)
	if os.IsNotExist(err) {
		return baseFileName, nil
	}
	dotIndex := strings.LastIndex(baseFileName, ".")
	nameParts := []string{baseFileName[:dotIndex], baseFileName[dotIndex:]}

	for i := 1; i < IterCount; i++ {
		resultName := fmt.Sprintf("%s(%d)%s", nameParts[0], i, nameParts[1])
		if _, err = os.Stat(path + "/" + resultName); os.IsNotExist(err) {
			return resultName, nil
		}
	}
	return "", ErrTypeExceeded
}
