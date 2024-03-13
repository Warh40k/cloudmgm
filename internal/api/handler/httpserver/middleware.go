package httpserver

import (
	"context"
	"github.com/Warh40k/cloud-manager/internal/api/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) CheckOwnership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "Middleware.CheckOwnership"
		log := h.log.With(
			slog.String("op", op),
		)
		userId, ok := r.Context().Value("user").(uuid.UUID)
		if !ok {
			log.With(slog.String("err", "failed to parse user id")).Error("failed to delete file")
			http.Error(w, "no user info", http.StatusForbidden)
			return
		}
		vmId, err := uuid.Parse(chi.URLParam(r, "volume_id"))
		if err != nil {
			log.With(slog.String("err", "failed to parse volume id")).Error("failed to delete file")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = h.services.CheckOwnership(userId, vmId); err != nil {
			log.With(slog.String("err", "failed to parse user id")).Error("failed to delete file")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		authHeader := r.Header.Get("Authorization")
		headSplit := strings.Split(authHeader, "Bearer ")
		if len(headSplit) == 2 {
			token = headSplit[1]
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		id, err := service.CheckJWT(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
