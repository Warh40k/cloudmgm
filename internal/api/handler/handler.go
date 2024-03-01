package handler

import (
	"github.com/Warh40k/cloud-manager/internal/api/service"
	genericMiddleware "github.com/Warh40k/cloud-manager/internal/middleware"
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
		r.Route("/volumes", func(r chi.Router) {
			r.Use(genericMiddleware.CheckAuth)
			r.Get("/", h.ListVolumes)
			r.Post("/", h.CreateVolume)
			r.Route("/{volume_id}", func(r chi.Router) {
				r.With(h.CheckOwnership).Get("/", h.GetVolume)
				r.With(h.CheckOwnership).Put("/", h.UpdateVolume)
				r.With(h.CheckOwnership).Delete("/", h.DeleteVolume)
				r.With(h.CheckOwnership).Post("/resize/", h.ResizeVolume)

				r.Route("/files", func(r chi.Router) {
					r.Post("/", h.UploadFile)
					r.Get("/{file_id}", h.GetFileInfo)
					r.Delete("/{file_id}", h.DeleteFile)
				})
			})
		})
	})

	return router
}
