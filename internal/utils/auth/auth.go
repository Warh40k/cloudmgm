package auth

import "golang.org/x/crypto/bcrypt"

func CheckAuth() bool {
	panic("not implemented")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(username, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func CheckJWT(token string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
