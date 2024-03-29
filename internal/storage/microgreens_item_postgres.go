package storage

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/riabkovK/microgreens/internal/domain"
)

type MicrogreensItemPostgres struct {
	db *sqlx.DB
}

func NewMicrogreensItemPostgres(db *sqlx.DB) *MicrogreensItemPostgres {
	return &MicrogreensItemPostgres{db: db}
}

func (mip *MicrogreensItemPostgres) Create(listId int, microgreensItem domain.MicrogreensItemRequest) (int, error) {
	tx, err := mip.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemID int
	createItemQuery := fmt.Sprintf(`INSERT INTO %s (name, description, price, microgreens_family_id) 
										   VALUES ($1, $2, $3, $4) RETURNING ID`, microgreensItemTable)

	row := tx.QueryRow(createItemQuery, microgreensItem.Name, microgreensItem.Description, microgreensItem.Price, microgreensItem.MicrogreensFamilyId)
	err = row.Scan(&itemID)
	if err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			return 0, err1
		}
		return 0, err
	}

	// Fill microgreens_list_items table
	createListItemQuery := fmt.Sprintf(`INSERT INTO %s (microgreens_list_id, microgreens_item_id) 
										   	   VALUES ($1, $2)`, microgreensListsItemsTable)

	_, err = tx.Exec(createListItemQuery, listId, itemID)
	if err != nil {
		if err1 := tx.Rollback(); err != nil {
			return 0, err1
		}
		return 0, err
	}

	// Fill microgreens_family_items table
	createFamilyItemQuery := fmt.Sprintf(`INSERT INTO %s (microgreens_family_id, microgreens_item_id)
										   	     VALUES ($1, $2)`, microgreensFamilyItemsTable)

	_, err = tx.Exec(createFamilyItemQuery, microgreensItem.MicrogreensFamilyId, itemID)
	if err != nil {
		if err1 := tx.Rollback(); err != nil {
			return 0, err1
		}
		return 0, err
	}

	return itemID, tx.Commit()
}

func (mip *MicrogreensItemPostgres) GetAll(userId, listId int) ([]domain.MicrogreensItem, error) {
	var items []domain.MicrogreensItem
	query := fmt.Sprintf(`SELECT mi.id, mi.name, mi.description, mi.price, mi.microgreens_family_id FROM %s AS mi 
                                INNER JOIN %s AS mli ON mli.microgreens_item_id = mi.id
                                INNER JOIN %s AS uml ON uml.microgreens_list_id = mli.microgreens_list_id
                                WHERE mli.microgreens_list_id = $1 AND uml.user_id = $2`,
		microgreensItemTable, microgreensListsItemsTable, usersMicrogreensListsTable)
	err := mip.db.Select(&items, query, listId, userId)

	return items, err
}

func (mip *MicrogreensItemPostgres) GetById(userId, itemId int) (domain.MicrogreensItem, error) {
	var item domain.MicrogreensItem
	query := fmt.Sprintf(`SELECT mi.id, mi.name, mi.description, mi.price, mi.microgreens_family_id FROM %s AS mi 
                                INNER JOIN %s AS mli ON mli.microgreens_item_id = mi.id
                                INNER JOIN %s AS uml ON uml.microgreens_list_id = mli.microgreens_list_id
                                WHERE mi.id = $1 AND uml.user_id = $2`,
		microgreensItemTable, microgreensListsItemsTable, usersMicrogreensListsTable)
	err := mip.db.Get(&item, query, itemId, userId)

	return item, err
}

func (mip *MicrogreensItemPostgres) Delete(userId, itemId int) (int, error) {
	query := fmt.Sprintf(`DELETE FROM %s AS mi USING %s AS mli, %s as uml
								WHERE mi.id = mli.microgreens_item_id AND mli.microgreens_list_id = uml.microgreens_list_id
								AND uml.user_id = $1 AND mi.id = $2`,
		microgreensItemTable, microgreensListsItemsTable, usersMicrogreensListsTable)
	result, err := mip.db.Exec(query, userId, itemId)
	if err != nil {
		return 0, err
	}

	// item is not exist if rowsAmount == 0
	rowsAmount, err := result.RowsAffected()

	return int(rowsAmount), err
}

func (mip *MicrogreensItemPostgres) Update(userId, itemId int, request domain.UpdateMicrogreensItemRequest) error {
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

	if request.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argID))
		args = append(args, *request.Price)
		argID++
	}

	if request.MicrogreensFamilyId != nil {
		setValues = append(setValues, fmt.Sprintf("microgreens_family_id=$%d", argID))
		args = append(args, *request.MicrogreensFamilyId)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s AS mi SET %s FROM %s AS mli, %s AS uml
                       			 WHERE mi.id = mli.microgreens_item_id AND mli.microgreens_list_id = uml.microgreens_list_id
								 AND uml.user_id = $%d AND mi.id = $%d`,
		microgreensItemTable, setQuery, microgreensListsItemsTable, usersMicrogreensListsTable, argID, argID+1)
	args = append(args, userId, itemId)

	_, err := mip.db.Exec(query, args...)
	return err
}
