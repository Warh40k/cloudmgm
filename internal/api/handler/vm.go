package handler

import (
	"encoding/json"
	"github.com/Warh40k/cloud-manager/internal/domain"
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
	//TODO implement
	panic("not implemented")

}

func (h *Handler) DeleteMachine(w http.ResponseWriter, r *http.Request) {
	//TODO implement
	panic("not implemented")

}

func (h *Handler) UpdateMachine(w http.ResponseWriter, r *http.Request) {
	//TODO implement
	panic("not implemented")

}
