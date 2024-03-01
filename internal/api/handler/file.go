package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	volumeId, err := uuid.Parse(chi.URLParam(r, "volume_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("uploaded_file")
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = UploadFile(volumeId, handler, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UploadFile(volumeId uuid.UUID, handler *multipart.FileHeader, file multipart.File) error {
	storagePath := viper.GetString("files.save_path") + "/" + volumeId.String()
	_, err := os.Stat(storagePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(storagePath, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}
	}

	dst, err := os.Create(storagePath + "/" + handler.Filename)
	defer dst.Close()
	if err != nil {
		return err
	}

	if _, err = io.Copy(dst, file); err != nil {
		return err
	}

	return nil
}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
