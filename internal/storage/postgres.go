package storage

import (
	"fmt"
	"github.com/riabkovK/microgreens/internal/config"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	usersTable             = "users"
	microgreensListTable   = "microgreens_list"
	microgreensFamilyTable = "microgreens_family"
	microgreensItemTable   = "microgreens_item"

	microgreensListsItemsTable  = "microgreens_list_items"
	usersMicrogreensListsTable  = "users_microgreens_lists"
	microgreensFamilyItemsTable = "microgreens_family_items"

	refreshSessionsTable = "refresh_sessions"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Username,
		cfg.Postgres.DBName,
		cfg.Postgres.Password,
		cfg.Postgres.SSLMode))

	if err != nil {
		logrus.WithError(err).Warning("preparing database connection")
		return nil, err
	}

	return db, err
}
