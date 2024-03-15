package internal

import "errors"

type MicrogreensList struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" validate:"required"`
	Description string `json:"description" db:"description"`
}

type MicrogreensFamily struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name" validate:"required"`
}

type MicrogreensItem struct {
	Id                  int    `json:"id" db:"id" `
	Name                string `json:"name" db:"name" validate:"required"`
	Description         string `json:"description" db:"description"`
	Price               int    `json:"price" db:"price" validate:"required"`
	MicrogreensFamilyId int    `json:"microgreens_family_id" db:"microgreens_family_id" validate:"required"`
}

type UsersMicrogreensList struct {
	Id                int
	UserId            int
	MicrogreensListId int
}

type MicrogreensListItems struct {
	Id                int
	MicrogreensListId int
	MicrogreensItemId int
}

type MicrogreensFamilyItems struct {
	Id                  int
	MicrogreensFamilyId int
	MicrogreensItemId   int
}

type UpdateMicrogreensListRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (receiver UpdateMicrogreensListRequest) Validate() error {
	if receiver.Name == nil && receiver.Description == nil {
		return errors.New("update microgreensList structure has no values")
	}
	return nil
}
