package domain

type User struct {
	Name     string `json:"name,omitempty"`
	Login    string `json:"login,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password,omitempty"`
}
