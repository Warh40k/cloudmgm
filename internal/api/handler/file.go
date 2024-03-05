package handler

import (
	"encoding/json"
	"github.com/Warh40k/cloud-manager/internal/api/handler/utils"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"os"
)

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	const op = "File.ListVolumeFiles"
	volumeIdParam := chi.URLParam(r, "volume_id")

	log := h.log.With(
		slog.String("op", op),
		slog.String("file id", volumeIdParam),
	)

	volumeId, err := uuid.Parse(volumeIdParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to parse multipart form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("uploaded_file")
	defer file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	const op = "File.Handler.DeleteFile"
	fileIdParam := chi.URLParam(r, "file_id")
	log := h.log.With(
		slog.String("op", op),
		slog.String("file id", fileIdParam),
	)
	fileId, err := uuid.Parse(fileIdParam)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to parse uuid")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.services.DeleteFile(fileId)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to delete file")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	const op = "File.Handler.GetFile"
	fileIdParam := chi.URLParam(r, "file_id")

	log := h.log.With(
		slog.String("op", op),
		slog.String("file id", fileIdParam),
	)

	fileId, err := uuid.Parse(fileIdParam)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to parse uuid")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := h.services.GetFile(fileId)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to get file")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(file)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to marshall json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) ListVolumeFiles(w http.ResponseWriter, r *http.Request) {
	const op = "File.Handler.ListVolumeFiles"
	volumeIdParam := chi.URLParam(r, "volume_id")

	log := h.log.With(
		slog.String("op", op),
		slog.String("volume id", volumeIdParam),
	)

	volumeId, err := uuid.Parse(volumeIdParam)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to marshall json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files, err := h.services.ListVolumeFiles(volumeId)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to list volumes")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(files)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to list volumes")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
