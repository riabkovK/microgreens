package storage

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/riabkovK/microgreens/internal/domain"
)

type MicrogreensFamilyPostgres struct {
	db *sqlx.DB
}

func NewMicrogreensFamilyPostgres(db *sqlx.DB) *MicrogreensFamilyPostgres {
	return &MicrogreensFamilyPostgres{db: db}
}

func (mfp *MicrogreensFamilyPostgres) Create(family domain.MicrogreensFamilyRequest) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, description) values ($1, $2) RETURNING id", microgreensFamilyTable)
	row := mfp.db.QueryRow(query, family.Name, family.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (mfp *MicrogreensFamilyPostgres) GetAll() ([]domain.MicrogreensFamily, error) {
	var families []domain.MicrogreensFamily
	query := fmt.Sprintf("SELECT * FROM %s", microgreensFamilyTable)
	err := mfp.db.Select(&families, query)

	return families, err
}

func (mfp *MicrogreensFamilyPostgres) GetById(familyId int) (domain.MicrogreensFamily, error) {
	var family domain.MicrogreensFamily
	query := fmt.Sprintf("SELECT * FROM %s AS mf WHERE mf.id = $1", microgreensFamilyTable)
	err := mfp.db.Get(&family, query, familyId)

	return family, err
}

func (mfp *MicrogreensFamilyPostgres) Delete(itemId int) (int, error) {
	query := fmt.Sprintf("DELETE FROM %s AS mf WHERE mf.id = $1", microgreensFamilyTable)

	result, err := mfp.db.Exec(query, itemId)
	if err != nil {
		return 0, err
	}

	// family is not exist if rowsAmount == 0
	rowsAmount, err := result.RowsAffected()

	return int(rowsAmount), err
}

func (mfp *MicrogreensFamilyPostgres) Update(familyId int, request domain.UpdateMicrogreensFamilyRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if request.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *request.Name)
		argID++
	}

	if request.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *request.Description)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s AS mf SET %s WHERE mf.id = $%d`,
		microgreensFamilyTable, setQuery, argID)
	args = append(args, familyId)

	_, err := mfp.db.Exec(query, args...)
	return err
}
