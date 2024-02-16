package handler

import (
	"encoding/json"
	"errors"
	"github.com/Warh40k/cloud-manager/internal/api/api_errors"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"io"
	"net/http"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.services.SignUp(user)
	if err != nil {
		if errors.Is(err, api_errors.ErrBadRequest) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
