package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/riabkovK/microgreens/internal"
)

type MicrogreensItemPostgres struct {
	db *sqlx.DB
}

func NewMicrogreensItemPostgres(db *sqlx.DB) *MicrogreensItemPostgres {
	return &MicrogreensItemPostgres{db: db}
}

func (mip *MicrogreensItemPostgres) Create(listId int, microgreensItem internal.MicrogreensItem) (int, error) {
	tx, err := mip.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemID int
	createItemQuery := fmt.Sprintf(`INSERT INTO %s (name, description, price, microgreens_family_id) 
										   VALUES ($1, $2, $3, $4, $5) RETURNING ID`, microgreensItemTable)

	row := tx.QueryRow(createItemQuery, microgreensItem.Name, microgreensItem.Description, microgreensItem.Price, microgreensItem.MicrogreensFamilyId)
	err = row.Scan(&itemID)
	if err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			return 0, err1
		}
		return 0, err
	}

	createListItemQuery := fmt.Sprintf(`INSERT INTO %s (microgreens_list_id, microgreens_item_id) 
										   	   VALUES ($1, $2)`, microgreensListsItemsTable)

	_, err = tx.Exec(createListItemQuery, listId, itemID)
	if err != nil {
		if err1 := tx.Rollback(); err != nil {
			return 0, err1
		}
		return 0, err
	}

	return itemID, tx.Commit()
}

func (mip *MicrogreensItemPostgres) GetAll(userId, listId int) ([]internal.MicrogreensItem, error) {
	var items []internal.MicrogreensItem
	query := fmt.Sprintf(`SELECT mi.id, mi.name, mi.description FROM %s AS mi 
                                INNER JOIN %s AS mli ON mli.microgreens_item_id = mi.id
                                INNER JOIN %s AS uml ON uml.microgreens_list_id = mli.microgreens_list_id
                                WHERE mli.microgreens_list_id = $1 AND uml.user_id = $2`,
		microgreensItemTable, microgreensListsItemsTable, usersMicrogreensListsTable)
	err := mip.db.Select(&items, query, listId, userId)

	return items, err
}

func (mip *MicrogreensItemPostgres) GetById(userId, itemId int) (internal.MicrogreensItem, error) {
	var item internal.MicrogreensItem
	query := fmt.Sprintf(`SELECT mi.id, mi.name, mi.description, mi.price, mi.microgreens_family_id FROM %s AS mi 
                                INNER JOIN %s AS mli ON mli.microgreens_item_id = mi.id
                                INNER JOIN %s AS uml ON uml.microgreens_list_id = mli.microgreens_list_id
                                WHERE mi.id = $1 AND uml.user_id = $2`,
		microgreensItemTable, microgreensListsItemsTable, usersMicrogreensListsTable)
	err := mip.db.Select(&item, query, itemId, userId)

	return item, err
}
