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
	createListQuery := fmt.Sprintf("INSERT INTO %s (name, micr) VALUES ($1, $2) RETURNING id", microgreensListTable)
	row := tx.QueryRow(createListQuery, list.Name, list.MicrogreensFamilyId)
}
