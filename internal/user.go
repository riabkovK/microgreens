package internal

type User struct {
	Id           int    `json:"-"`
	Email        string `json:"email" validate:"required,email"`
	Name         string `json:"name" validate:"required"`
	Username     string `json:"username" validate:"required,max=25"`
	PasswordHash string `json:"password" validate:"required,sha256"`
}
