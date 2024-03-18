package internal

import "errors"

// Default structures

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

// Structures for updating

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

type UpdateMicrogreensItemRequest struct {
	Name                *string `json:"name"`
	Description         *string `json:"description"`
	Price               *int    `json:"price" validate:"min=0"`
	MicrogreensFamilyId *int    `json:"microgreens_family_id" validate:"min=0"`
}

func (receiver UpdateMicrogreensItemRequest) Validate() error {
	if receiver.Name == nil && receiver.Description == nil && receiver.Price == nil && receiver.MicrogreensFamilyId == nil {
		return errors.New("update microgreensItem structure has no values")
	}
	return nil
}
