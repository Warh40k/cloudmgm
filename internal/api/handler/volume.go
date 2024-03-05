package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/api/handler/utils"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func (h *Handler) ListVolumes(w http.ResponseWriter, r *http.Request) {
	const op = "Volume.Handler.ListVolumes"
	log := h.log.With(
		slog.String("op", op),
	)

	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		log.Error("failed to delete file")
		http.Error(w, "wrong user id", http.StatusBadRequest)
		return
	}

	vms, err := h.services.ListVolume(userId)
	if err != nil {
		log.With(
			slog.String("err", err.Error()),
			slog.String("user id", userId.String()),
		).Error("failed to delete file")

		w.WriteHeader(http.StatusNotFound)
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}
	if vms == nil {
		vms = []domain.Volume{}
	}

	responseText, err := json.Marshal(vms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseText)
}

func (h *Handler) GetVolume(w http.ResponseWriter, r *http.Request) {
	vmId, err := uuid.Parse(chi.URLParam(r, "volume_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vm, err := h.services.GetVolume(vmId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(vm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *Handler) CreateVolume(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(
		slog.String("op", "Middleware.CheckOwnership"),
	)
	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		log.With(slog.String("err", "failed to parse user id")).Error("failed to delete file")

		http.Error(w, "wrong user id", http.StatusBadRequest)
		return
	}

	var vm domain.Volume

	err := json.NewDecoder(r.Body).Decode(&vm)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to delete file")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(vm)
	if err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		log.With(slog.String("err", errs.Error())).Error("failed to validate input volume")
		http.Error(w, fmt.Sprintf("Validation error: %s", errs), http.StatusBadRequest)
		return
	}
	_, err = h.services.CreateVolume(userId, vm)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to create volume")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteVolume(w http.ResponseWriter, r *http.Request) {
	const op = "Volume.Handler.DeleteVolume"
	log := h.log.With(
		slog.String("op", op),
	)

	vmId, err := uuid.Parse(chi.URLParam(r, "volume_id"))
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to get volume id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.services.DeleteVolume(vmId)
	if err != nil {
		log.With(slog.String("err", err.Error())).Error("failed to delete volume")
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (h *Handler) UpdateVolume(w http.ResponseWriter, r *http.Request) {
	vmId, err := uuid.Parse(chi.URLParam(r, "volume_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var vm domain.Volume
	err = json.NewDecoder(r.Body).Decode(&vm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vm.Id = vmId
	err = h.services.UpdateVolume(vm)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

type ResizeVmIn struct {
	Increase bool   `json:"increase" validate:"required"`
	Amount   string `json:"amount" validate:"required"`
}

func (h *Handler) ResizeVolume(w http.ResponseWriter, r *http.Request) {
	vmId, err := uuid.Parse(chi.URLParam(r, "volume_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var input ResizeVmIn

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size, err := utils.ConvertSizeToBytes(input.Amount, input.Increase)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.services.ResizeVolume(vmId, size)
}
