package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/riabkovK/microgreens/internal"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (receiver *AuthPostgres) CreateUser(user internal.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, name, username, password_hash) values ($1, $2, $3, $4) RETURNING id", usersTable)
	row := receiver.db.QueryRow(query, user.Email, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (receiver *AuthPostgres) GetUser(email, password string) (internal.User, error) {
	var user internal.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := receiver.db.Get(&user, query, email, password)

	return user, err
}
