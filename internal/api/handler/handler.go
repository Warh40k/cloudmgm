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
		r.Route("/machines", func(r chi.Router) {
			r.Use(middleware2.CheckAuth)
			r.Get("/", h.ListMachines)
			r.Get("/{machine_id}", h.GetMachine)
			r.Post("/", h.CreateMachine)
			r.Patch("/{machine_id}", h.UpdateMachine)
			r.Delete("/{machine_id}", h.DeleteMachine)
		})
	})

	return router
}
