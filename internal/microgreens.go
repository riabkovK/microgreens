package internal

import "errors"

type MicrogreensList struct {
	Id                  int    `json:"id" db:"id"`
	Name                string `json:"name" db:"name" validate:"required"`
	Description         string `json:"description" db:"description"`
	MicrogreensFamilyId int    `json:"microgreens_family_id" db:"microgreens_family_id"`
}

type MicrogreensFamily struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name" validate:"required"`
}

type MicrogreensItem struct {
	Id                  int    `json:"id" db:"id" `
	Name                string `json:"name" db:"name" validate:"required"`
	Description         string `json:"description" db:"description"`
	MicrogreensFamilyId int    `json:"microgreens_family_id" db:"microgreens_family_id"`
	Price               int    `json:"price" db:"price" validate:"required"`
}

type UsersMicrogreensList struct {
	Id                int
	UserId            int
	MicrogreensListId int
}

type UpdateMicrogreensListRequest struct {
	Name                *string `json:"name"`
	Description         *string `json:"description"`
	MicrogreensFamilyId *int    `json:"microgreens_family_id"`
}

func (receiver UpdateMicrogreensListRequest) Validate() error {
	if receiver.Name == nil && receiver.Description == nil && receiver.MicrogreensFamilyId == nil {
		return errors.New("update microgreensList structure has no values")
	}
	return nil
}
