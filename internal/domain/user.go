package domain

import "time"

type User struct {
	Id           int       `json:"-" db:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
	LastVisitAt  time.Time `json:"last_visit_at"`
}
