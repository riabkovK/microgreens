package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/riabkovK/microgreens/internal"
	"strings"
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

func (mlsp *MicrogreensListPostgres) GetAll(userId int) ([]internal.MicrogreensList, error) {
	var lists []internal.MicrogreensList
	query := fmt.Sprintf(`SELECT ml.id, ml.name, ml.description FROM %s AS ml 
                                INNER JOIN %s as uml ON ml.id = uml.microgreens_list_id 
                                WHERE uml.user_id = $1`,
		microgreensListTable, usersMicrogreensListsTable)
	err := mlsp.db.Select(&lists, query, userId)

	return lists, err
}

func (mlsp *MicrogreensListPostgres) GetById(userId, listId int) (internal.MicrogreensList, error) {
	list := internal.MicrogreensList{}
	query := fmt.Sprintf(`SELECT ml.id, ml.name, ml.description FROM %s AS ml 
                                INNER JOIN %s as uml ON ml.id = uml.microgreens_list_id 
                                WHERE uml.user_id = $1 AND uml.microgreens_list_id = $2`,
		microgreensListTable, usersMicrogreensListsTable)
	err := mlsp.db.Get(&list, query, userId, listId)

	return list, err
}

func (mlsp *MicrogreensListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf(`DELETE FROM %s AS ml USING %s AS uml
								WHERE ml.id = uml.microgreens_list_id AND uml.user_id = $1 AND uml.microgreens_list_id = $2`,
		microgreensListTable, usersMicrogreensListsTable)
	_, err := mlsp.db.Exec(query, userId, listId)

	return err
}

func (mlsp *MicrogreensListPostgres) Update(userId, listId int, request internal.UpdateMicrogreensListRequest) error {
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

	if request.MicrogreensFamilyId != nil {
		setValues = append(setValues, fmt.Sprintf("microgreens_family_id=$%d", argID))
		args = append(args, *request.MicrogreensFamilyId)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s AS ml SET %s FROM %s AS uml
								 WHERE ml.id = uml.microgreens_list_id AND uml.microgreens_list_id=$%d 
									AND uml.user_id=$%d`,
		microgreensListTable, setQuery, usersMicrogreensListsTable, argID, argID+1)
	args = append(args, listId, userId)

	_, err := mlsp.db.Exec(query, args...)
	return err
}
