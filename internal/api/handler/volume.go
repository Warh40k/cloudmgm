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
	"net/http"
)

func (h *Handler) ListVolumes(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vms, err := h.services.ListVolume(userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if vms == nil {
		vms = []domain.Volume{}
	}

	responseText, err := json.Marshal(vms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseText)
}

func (h *Handler) GetVolume(w http.ResponseWriter, r *http.Request) {
	vmId, err := uuid.Parse(chi.URLParam(r, "volume_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vm, err := h.services.GetVolume(vmId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(vm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *Handler) CreateVolume(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var vm domain.Volume

	err := json.NewDecoder(r.Body).Decode(&vm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(vm)
	if err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		http.Error(w, fmt.Sprintf("Validation error: %s", errs), http.StatusBadRequest)
	}
	_, err = h.services.CreateVolume(userId, vm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Handler) DeleteVolume(w http.ResponseWriter, r *http.Request) {
	vmId, err := uuid.Parse(chi.URLParam(r, "volume_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.services.DeleteVolume(vmId)
	if err != nil {
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
