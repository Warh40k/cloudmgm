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
			http.Error(w, "no user info", http.StatusForbidden)
			return
		}

		vmId, err := uuid.Parse(chi.URLParam(r, "volume_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = h.services.CheckOwnership(userId, vmId); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
