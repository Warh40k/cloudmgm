package handler

import (
	"github.com/Warh40k/cloud-manager/internal/api/service"
	middleware2 "github.com/Warh40k/cloud-manager/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/auth", h.SignIn)
		r.Post("/register", h.SignUp)
		r.With(middleware2.CheckAuth).Get("/pong", h.Pong)
	})

	return router
}
