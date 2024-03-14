package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/riabkovK/microgreens/internal"
)

type MicrogreensListPostgres struct {
	db *sqlx.DB
}

func NewMicrogreensListPostgres(db *sqlx.DB) *MicrogreensListPostgres {
	return &MicrogreensListPostgres{db: db}
}

func (mlsp *MicrogreensListPostgres) Create(userId int, list internal.MicrogreensList) (int, error) {
	tx, err := mlsp.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createMicrogreensListQuery := fmt.Sprintf("INSERT INTO %s (name, description) VALUES ($1, $2) RETURNING id", microgreensListTable)
	row := tx.QueryRow(createMicrogreensListQuery, list.Name, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersMicrogreensListQuery := fmt.Sprintf("INSERT INTO %s (user_id, microgreens_list_id) VALUES ($1, $2)", usersMicrogreensListsTable)
	_, err = tx.Exec(createUsersMicrogreensListQuery, userId, id)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}
