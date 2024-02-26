package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) CheckOwnership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value("user").(uuid.UUID)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		vmId, err := uuid.Parse(chi.URLParam(r, "machine_id"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = h.services.CheckOwnership(userId, vmId); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}