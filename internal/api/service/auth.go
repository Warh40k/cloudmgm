package service

import (
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"time"
)

const (
	salt       = "fjlsj2374slfjsd728vvnts"
	tokenTTL   = 12 * time.Hour
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFsdsfdjH"
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (a AuthService) SignUp(user domain.User) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (a AuthService) SignIn(username, password string) (int, error) {
	//TODO implement me
	panic("implement me")
}
