package handler

import (
	"encoding/json"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"io"
	"log"
	"net/http"
)

var users []domain.User

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
	users = append(users, user)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(users)
	if err != nil {
		log.Fatalf("Jopa %s", err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {

}
