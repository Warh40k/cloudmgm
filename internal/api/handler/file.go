package handler

import (
	"github.com/Warh40k/cloud-manager/internal/api/handler/utils"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
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

	dirPath := viper.GetString("files.save_path") + "/" + volumeId.String()
	fileName, err := utils.GetFileName(dirPath, handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModeDir|os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	dst, err := os.Create(dirPath + "/" + fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	err = utils.UploadFile(dirPath, file, dst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var fileEntity = domain.File{
		VolumeId: volumeId,
		Name:     fileName,
		Size:     handler.Size,
		Link:     "",
	}

	_, err = h.services.CreateFile(fileEntity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetFileInfo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
