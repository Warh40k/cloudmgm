package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)

// TODO Реализовать методы
func (h *Handler) ListMachines(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vms, err := h.services.ListVm(userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if vms == nil {
		vms = []domain.VirtualMachine{}
	}

	responseText, err := json.Marshal(vms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseText)
}

func (h *Handler) GetMachine(w http.ResponseWriter, r *http.Request) {
	vmId, err := uuid.Parse(chi.URLParam(r, "machine_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vm, err := h.services.GetVm(vmId)
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

func (h *Handler) CreateMachine(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var vm domain.VirtualMachine

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
	_, err = h.services.CreateVm(userId, vm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Handler) DeleteMachine(w http.ResponseWriter, r *http.Request) {
	vmId, err := uuid.Parse(chi.URLParam(r, "machine_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.services.DeleteVm(vmId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (h *Handler) UpdateMachine(w http.ResponseWriter, r *http.Request) {
	vmId, err := uuid.Parse(chi.URLParam(r, "machine_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var vm domain.VirtualMachine
	err = json.NewDecoder(r.Body).Decode(&vm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vm.Id = vmId
	err = h.services.UpdateVm(vm)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
