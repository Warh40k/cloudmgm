package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"sync"
)

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	const op = "File.UploadFile"
	volumeIdParam := chi.URLParam(r, "volume_id")

	log := h.log.With(
		slog.String("op", op),
		slog.String("file id", volumeIdParam),
	)

	volumeId, err := uuid.Parse(volumeIdParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to parse multipart form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("uploaded_file")
	defer file.Close()
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to upload file")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fs := afero.NewOsFs()
	volumePath := viper.GetString("files.save_path") + "/" + volumeId.String()
	fileName, err := h.services.UploadFile(volumePath, file, header.Filename, fs)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to save file")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// save file info in db
	var fileEntity = domain.File{
		VolumeId: volumeId,
		Name:     fileName,
		Size:     header.Size,
		Link:     "",
	}

	_, err = h.services.CreateFile(fileEntity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UploadError struct {
	Filename string
	Message  string
	Err      error
}

func (e UploadError) Error() string {
	return fmt.Sprintf("%s: %s", e.Filename, e.Message)
}

func (h *Handler) UploadMultipleFiles(w http.ResponseWriter, r *http.Request) {
	const op = "File.UploadFile"
	volumeIdParam := chi.URLParam(r, "volume_id")

	log := h.log.With(
		slog.String("op", op),
		slog.String("file id", volumeIdParam),
	)

	volumeId, err := uuid.Parse(volumeIdParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to parse multipart form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	errs := make(chan error, len(files))
	var results = make([]uuid.UUID, len(files))

	wg := sync.WaitGroup{}
	for i, fileHeader := range files {
		wg.Add(1)
		go h.UploadFileWorker(i, fileHeader, &wg, errs, volumeId, results)
	}
	wg.Wait()
	close(errs)

	for err = range errs {
		var uploadErr UploadError
		errors.As(err, &uploadErr)
		log.With(slog.String("err", err.Error()),
			slog.String("filename", uploadErr.Filename)).
			Error(uploadErr.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp, err := json.Marshal(results)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to create json response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *Handler) UploadFileWorker(i int, fileHeader *multipart.FileHeader, wg *sync.WaitGroup, errs chan<- error, volumeId uuid.UUID, results []uuid.UUID) {
	defer wg.Done()

	file, err := fileHeader.Open()
	if err != nil {
		errs <- UploadError{
			Filename: fileHeader.Filename,
			Message:  "failed to open file",
			Err:      err,
		}
		return
	}
	defer file.Close()
	fs := afero.NewOsFs()
	volumePath := viper.GetString("files.save_path") + "/" + volumeId.String()
	saveName, err := h.services.UploadFile(volumePath, file, fileHeader.Filename, fs)
	if err != nil {
		errs <- UploadError{
			Filename: fileHeader.Filename,
			Message:  "failed to upload file",
			Err:      err,
		}
		return
	}
	id, err := h.services.CreateFile(domain.File{
		VolumeId: volumeId,
		Name:     saveName,
		Size:     fileHeader.Size,
	})
	if err != nil {
		errs <- UploadError{
			Filename: fileHeader.Filename,
			Message:  "failed to save file info",
			Err:      err,
		}
		return
	}
	results[i] = id
}

func (h *Handler) DownloadMultipleFiles(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	const op = "File.Handler.DeleteFileInfo"
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
		log.With(slog.String("err", err.Error())).Error("failed to delete file info")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

}

func (h *Handler) GetFileInfo(w http.ResponseWriter, r *http.Request) {
	const op = "File.Handler.GetFileInfo"
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

	file, err := h.services.GetFileInfo(fileId)
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
	const op = "File.Handler.DownloadFile"
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

	/*file, err := h.services.GetFileInfo(fileId)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to get file")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()*/

	fileInfo, err := h.services.GetFileInfo(fileId)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to get file")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	attHeader := mime.FormatMediaType("attachment", map[string]string{"filename": fileInfo.Name})
	r.Header.Set("Content-Disposition", attHeader)
	r.Header.Set("Content-Type", "application/octet-stream")

	path := fileInfo.GetPath()
	http.ServeFile(w, r, path)
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
