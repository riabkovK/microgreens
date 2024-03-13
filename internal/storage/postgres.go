package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	usersTable                 = "users"
	microgreensListTable       = "microgreens_list"
	microgreensFamilyTable     = "microgreens_family"
	microgreensItemTable       = "microgreens_item"
	usersMicrogreensListsTable = "users_microgreens_lists"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DBName,
		cfg.Password,
		cfg.SSLMode))

	if err != nil {
		logrus.WithError(err).Warning("preparing database connection")
		return nil, err
	}

	return db, err
}
