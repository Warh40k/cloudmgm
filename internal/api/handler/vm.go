package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)

// TODO Реализовать методы
func (h *Handler) ListMachines(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vms, err := h.services.ListVm(id)
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
	//TODO implement
	panic("not implemented")
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
	err = h.services.CreateVm(userId, vm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Handler) DeleteMachine(w http.ResponseWriter, r *http.Request) {
	//TODO implement
	panic("not implemented")

}

func (h *Handler) UpdateMachine(w http.ResponseWriter, r *http.Request) {
	//TODO implement
	panic("not implemented")

}
