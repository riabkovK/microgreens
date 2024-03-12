package internal

type User struct {
	Id           int    `json:"-"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}
